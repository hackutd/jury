#[rocket::get("/")]
pub fn home() -> &'static str {
    "Hello world!"
}
