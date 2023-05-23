use mongodb::Database;
use rocket::response::stream::{Event, EventStream};
use rocket::{http::Status, serde::json::Json, State};
use std::env;
use std::sync::Arc;
use tokio::sync::mpsc::{channel, Sender};
use tokio::sync::Mutex;

use crate::util::tasks::NewSender;
use crate::util::types::ClockState;
use crate::{
    db::admin::aggregate_stats,
    util::types::{AdminLogin, Stats},
};

use super::util::AdminPassword;

#[rocket::post("/admin/login", data = "<body>")]
pub async fn login(body: Json<AdminLogin<'_>>) -> (Status, String) {
    // Get correct passsword from environmental variables
    let correct = env::var("JURY_ADMIN_PASSWORD").expect("JURY_ADMIN_PASSWORD not defined");

    // Return status based on if the password is correct or not
    if correct == body.password.to_string() {
        (Status::Ok, correct)
    } else {
        (
            Status::BadRequest,
            "Invalid or missing password field".to_string(),
        )
    }
}

#[rocket::get("/admin/stats")]
pub async fn get_stats(
    _password: AdminPassword,
    db: &State<Arc<Database>>,
) -> (Status, Json<Stats>) {
    match aggregate_stats(&db).await {
        Ok(stats) => (Status::Ok, Json(stats)),
        Err(_) => (Status::InternalServerError, Json(Stats::default())),
    }
}

#[rocket::get("/admin/sync")]
pub async fn req_sync(
    _password: AdminPassword,
    tx_new_sender: &State<Sender<NewSender>>,
) -> (Status, Option<EventStream![Event + '_]>) {
    // Create new message channel for the current admin
    let (tx, mut rx) = channel::<String>(10);

    // Send message sender to new_sender message channel
    match tx_new_sender.send(NewSender { connect: true, tx }).await {
        Ok(()) => (),
        Err(e) => {
            eprintln!("Unable to send message to new sender message channel: {e}",);
            return (Status::InternalServerError, None);
        }
    };

    // If successfully sent message sender, listen for messages on that channel
    // 2 events possible: "stats" (update on stats) and "clock" (update on clock)
    (
        Status::Ok,
        Some(EventStream! {
            loop {
                // Wait for message from channel
                let msg = match &rx.recv().await {
                    Some(x) => x.clone(),
                    None => continue,
                };
                yield Event::data(msg);
            }
        }),
    )
}

#[rocket::get("/admin/clock")]
pub async fn clock(
    _password: AdminPassword,
    clock_state: &State<Arc<Mutex<ClockState>>>,
) -> (Status, Json<ClockState>) {
    // Get clock state
    let state = clock_state.lock().await;

    // Return clock state
    (Status::Ok, Json(state.get_copy()))
}

#[rocket::post("/admin/clock/pause")]
pub async fn pause_clock(
    _password: AdminPassword,
    clock_tx: &State<Sender<String>>,
) -> (Status, String) {
    clock_message(clock_tx, "pause".to_string()).await
}

#[rocket::post("/admin/clock/unpause")]
pub async fn unpause_clock(
    _password: AdminPassword,
    clock_tx: &State<Sender<String>>,
) -> (Status, String) {
    clock_message(clock_tx, "unpause".to_string()).await
}

#[rocket::post("/admin/clock/reset")]
pub async fn reset_clock(
    _password: AdminPassword,
    clock_tx: &State<Sender<String>>,
) -> (Status, String) {
    clock_message(clock_tx, "reset".to_string()).await
}

async fn clock_message(clock_tx: &State<Sender<String>>, msg: String) -> (Status, String) {
    match clock_tx.send(msg.clone()).await {
        Ok(()) => (Status::Ok, format!("Clock {} successful", msg)),
        Err(e) => (
            Status::InternalServerError,
            format!("Unable to {} clock: {}", msg, e),
        ),
    }
}
