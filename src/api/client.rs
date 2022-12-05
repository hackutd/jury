use rocket::response::Redirect;
use rocket_dyn_templates::{context, Template};

use std::env;

use super::guards::{AdminPassword, Token};

#[rocket::get("/")]
pub fn home() -> Template {
    Template::render("index", context! { title: get_title() })
}

#[rocket::get("/judge/login")]
pub fn judge_login() -> Result<Template, Redirect> {
    Ok(Template::render(
        "judge-login",
        context! { title: get_title() },
    ))
}

#[rocket::get("/judge")]
pub fn judge(_token: Token) -> Template {
    Template::render("judge", context! { title: get_title() })
}

#[rocket::get("/judge/welcome")]
pub fn judge_welcome(_token: Token) -> Template {
    Template::render("judge-welcome", context! { title: get_title() })
}

#[rocket::get("/admin/login")]
pub fn admin_login() -> Template {
    Template::render("admin-login", context! { title: get_title() })
}

#[rocket::get("/admin")]
pub fn admin(_password: AdminPassword) -> Template {
    Template::render("admin", context! { title: get_title() })
}

fn get_title() -> String {
    env::var("GAVEL_NAME").expect("GAVEL_NAME not defined")
}
