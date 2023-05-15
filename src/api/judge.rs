use std::sync::Arc;

use bson::doc;
use mongodb::Database;
use rand::{distributions::Alphanumeric, Rng};
use rocket::form::Form;
use rocket::http::Status;
use rocket::serde::json::Json;
use rocket::State;

use crate::util::types::CsvUpload;
use crate::{
    db::judge::{
        find_judge_by_code, insert_judge, insert_judges, read_welcome, update_judge_token,
    },
    try_500,
    util::parse_csv::parse_judge_csv,
};

use super::{
    request_types::{Login, NewJudge},
    util::{parse_csv_from_data, AdminPassword, Token},
};

#[rocket::post("/judge/login", data = "<body>")]
pub async fn login(db: &State<Arc<Database>>, body: Json<Login<'_>>) -> (Status, String) {
    // Find judge from db using code
    let judge = try_500!(
        find_judge_by_code(db, body.code).await,
        "Unable to process or find code"
    );

    // Generate random 16-character token
    let token: String = rand::thread_rng()
        .sample_iter(&Alphanumeric)
        .take(16)
        .map(char::from)
        .collect();

    // Update the token in the database and return token
    match update_judge_token(db, &judge.code, &token).await {
        Ok(_) => (Status::Ok, (format!("{}", token))),
        Err(_) => (
            Status::InternalServerError,
            "Unable to create token".to_string(),
        ),
    }
}

#[rocket::post("/judge/new", data = "<body>")]
pub async fn new_judge(
    db: &State<Arc<Database>>,
    body: Json<NewJudge<'_>>,
    _password: AdminPassword,
) -> (Status, String) {
    match insert_judge(db, body.0).await {
        Ok(_) => (Status::Accepted, "{}".to_string()),
        Err(status) => (status, "Invalid code".to_string()),
    }
}

#[rocket::post("/judge/csv", data = "<upload>")]
pub async fn add_judges_csv(
    upload: Form<CsvUpload>,
    db: &State<Arc<Database>>,
    _password: AdminPassword,
) -> (Status, String) {
    // Parse the CSV data
    let judges = try_500!(
        parse_judge_csv(upload.csv.clone(), upload.has_header).await,
        "Unable to parse CSV"
    );

    // If there are no judges, return Ok
    let num_judges = judges.len();
    if num_judges == 0 {
        return (Status::Ok, "0 judges added".to_string());
    }

    // Save the parsed CSV data to the database
    try_500!(
        insert_judges(db, judges).await,
        "Unable to insert judges into database"
    );

    (
        Status::Ok,
        format!("{} judges added", num_judges).to_string(),
    )
}

#[rocket::post("/judge/welcome")]
pub async fn judge_read_welcome(db: &State<Arc<Database>>, token: Token) -> (Status, String) {
    match read_welcome(db, &token.0).await {
        Ok(_) => (Status::Accepted, "{}".to_string()),
        Err(_) => (
            Status::InternalServerError,
            "Internal Server Error".to_string(),
        ),
    }
}
