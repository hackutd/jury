use std::sync::Arc;

use crate::db::project::insert_projects;
use crate::util::parse_csv::devpost_integration;
use mongodb::Database;
use rocket::data::Data;
use rocket::http::Status;
use rocket::State;

use super::util::parse_csv_from_data;

#[rocket::post("/project/csv", data = "<csv>")]
pub async fn add_devpost_csv(csv: Data<'_>, db: &State<Arc<Database>>) -> (Status, String) {
    // TODO: Add admin token check
    let cut_str = parse_csv_from_data(csv).await;

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
