use std::sync::Arc;

use bson::doc;
use mongodb::Database;
use rand::{distributions::Alphanumeric, Rng};
use rocket::data::Data;
use rocket::http::Status;
use rocket::serde::json::Json;
use rocket::State;

use crate::{db::judge::{find_judge_by_code, insert_judge, read_welcome, update_judge_token, insert_judges}, util::parse_csv::parse_judge_csv};

use super::{
    request_types::{Login, NewJudge},
    util::{parse_csv_from_data, AdminPassword, Token},
};

#[rocket::post("/judge/login", data = "<body>")]
pub async fn login(db: &State<Arc<Database>>, body: Json<Login<'_>>) -> (Status, String) {
    // Find judge from db using code
    let judge = match find_judge_by_code(db, body.code).await {
        Ok(j) => j,
        Err(status) => {
            return (status, "Unable to process or find code".to_string());
        }
    };

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

#[rocket::post("/judge/csv", data="<csv>")]
pub async fn add_judges_csv(csv: Data<'_>, db: &State<Arc<Database>>) -> (Status, String) {
    // TODO: Add admin token check
    let cut_str = parse_csv_from_data(csv).await;

    // Parse the CSV data
    let judges = match parse_judge_csv(cut_str).await {
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
    match insert_judges(db, judges).await {
        Ok(_) => (),
        Err(e) => {
            eprintln!("Unable to insert judges into database: {e}",);
            return (
                Status::InternalServerError,
                format!("Unable to insert judges into database: {e}").to_string(),
            );
        }
    };

    (Status::Ok, "".to_string())
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
