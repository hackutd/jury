use std::sync::Arc;

use crate::db::project::insert_projects;
use crate::try_status;
use crate::util::parse_csv::devpost_integration;
use mongodb::Database;
use rocket::data::Data;
use rocket::http::Status;
use rocket::State;

use super::util::{parse_csv_from_data, AdminPassword};

#[rocket::post("/project/devpost", data = "<csv>")]
pub async fn add_devpost_csv(
    csv: Data<'_>,
    db: &State<Arc<Database>>,
    _password: AdminPassword,
) -> (Status, String) {
    let cut_str = try_status!(
        parse_csv_from_data(csv).await,
        "Unable to parse CSV from data",
        Status::BadRequest
    );

    // Parse the CSV data
    let parsed_csv = try_status!(
        devpost_integration(cut_str, db).await,
        "Unable to parse CSV",
        Status::BadRequest
    );

    // Save the parsed CSV data to the database
    try_status!(
        insert_projects(&db, parsed_csv).await,
        "Unable to insert projects into database",
        Status::InternalServerError
    );

    (Status::Ok, "".to_string())
}
