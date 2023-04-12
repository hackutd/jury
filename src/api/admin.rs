use mongodb::Database;
use rocket::data::{Data, ToByteUnit};
use rocket::response::stream::{Event, EventStream};
use rocket::{http::Status, serde::json::Json, State};
use std::env;
use std::sync::Arc;
use tokio::sync::mpsc::{channel, Sender};

use crate::db::admin::insert_projects;
use crate::util::parse_csv::devpost_integration;
use crate::util::tasks::NewSender;
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

const CSV_HEADER: &str = "Content-Type: text/csv\r\n\r\n";
const CSV_FOOTER: &str = "\r\n----------------------------";

#[rocket::post("/admin/csv", data = "<csv>")]
pub async fn add_projects_csv(csv: Data<'_>, db: &State<Arc<Database>>) -> (Status, String) {
    // TODO: Add admin token check

    // Extract all data and store into a string
    // TODO: No unwrap
    let str = csv.open(16.mebibytes()).into_string().await.unwrap();

    // Find the start and end of the CSV data
    // and cut out the CSV data from the raw string
    let start = str.find(CSV_HEADER).unwrap() + CSV_HEADER.len();
    let end = str[start..].find(CSV_FOOTER).unwrap();
    let cut_str = &str[start..(start + end)];

    // Parse the CSV data
    let parsed_csv = match devpost_integration(cut_str, db).await {
        Ok(p) => p,
        Err(e) => {
            eprintln!("Unable to parse CSV: {e}",);
            return (
                Status::InternalServerError,
                format!("Unable to parse CSV: {e}").to_string(),
            );
        }
    };

    // Save the parsed CSV data to the database
    match insert_projects(&db, parsed_csv).await {
        Ok(_) => (),
        Err(e) => {
            eprintln!("Unable to insert projects into database: {e}",);
            return (
                Status::InternalServerError,
                format!("Unable to insert projects into database: {e}").to_string(),
            );
        }
    };

    (Status::Ok, "".to_string())
}
