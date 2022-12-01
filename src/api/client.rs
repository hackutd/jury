use rocket_dyn_templates::{context, Template};

use std::env;

#[rocket::get("/")]
pub fn home() -> Template {
    Template::render(
        "index",
        context! { title: env::var("GAVEL_NAME").expect("GAVEL_NAME not defined") },
    )
}
