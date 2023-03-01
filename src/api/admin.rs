use mongodb::Database;
use rocket::{http::Status, serde::json::Json, State};
use std::env;

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
pub async fn get_stats(_password: AdminPassword, db: &State<Database>) -> (Status, Json<Stats>) {
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
