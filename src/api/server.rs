use rocket::response::status;
use rocket::serde::{json::Json, Deserialize};

#[derive(Deserialize)]
#[serde(crate = "rocket::serde")]
pub struct Login<'r> {
    code: &'r str,
}

#[rocket::post("/judge/login", data = "<body>")]
pub fn judge_login(body: Json<Login<'_>>) -> status::Accepted<String> {
    status::Accepted(Some(format!("{}", body.code)))
}
