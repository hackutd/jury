use rocket::{
    http::Status,
    serde::{json::Json, Deserialize},
};
use std::env;

#[derive(Deserialize)]
#[serde()]
pub struct AdminLogin<'r> {
    password: &'r str,
}

#[rocket::post("/admin/login", data = "<body>")]
pub async fn login(body: Json<AdminLogin<'_>>) -> (Status, String) {
    // Get correct passsword from environmental variables
    let correct = env::var("GAVEL_ADMIN_PASSWORD").expect("GAVEL_ADMIN_PASSWORD not defined");

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
