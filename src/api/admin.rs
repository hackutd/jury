use mongodb::Database;
use rocket::response::stream::{Event, EventStream};
use rocket::{http::Status, serde::json::Json, State};
use std::env;
use std::sync::Arc;
use tokio::sync::mpsc::{channel, Sender};

use crate::{
    db::admin::aggregate_stats,
    util::types::{AdminLogin, Stats},
};

use super::{tasks::NewSender, util::AdminPassword};

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
pub async fn get_stats(_password: AdminPassword, db: &State<Arc<Database>>) -> (Status, Json<Stats>) {
    match aggregate_stats(&db).await {
        Ok(stats) => (Status::Ok, Json(stats)),
        Err(_) => (
            Status::InternalServerError,
            Json(Stats {
                projects: 0,
                seen: 0,
                votes: 0,
                time: 0,
                avg_mu: 0.0,
                avg_sigma: 0.0,
                judges: 0,
            }),
        ),
    }
}

#[rocket::get("/admin/sync")]
pub async fn req_sync(
    _password: AdminPassword,
    rx_new_sender: &State<Sender<NewSender>>,
) -> (Status, Option<EventStream![Event + '_]>) {
    // Create new message channel for the current admin
    let (tx, mut rx) = channel::<String>(10);

    // Send message sender to new_sender message channel
    match rx_new_sender.send(NewSender { connect: true, tx }).await {
        Ok(()) => (),
        Err(e) => {
            eprintln!("Unable to send message to new sender message channel: {e}",);
            return (Status::InternalServerError, None);
        }
    };

    // If successfully sent message sender, listen for messages on that channel
    (
        Status::Ok,
        Some(EventStream! {
            loop {
                // Wait for message from channel
                match &rx.recv().await {
                    Some(x) => x,
                    None => continue,
                };
                yield Event::data("update");
            }
        }),
    )
}
