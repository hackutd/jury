use rocket_dyn_templates::{context, Template};

use std::env;

#[rocket::get("/")]
pub fn home() -> Template {
    Template::render("index", context! { title: get_title() })
}

#[rocket::get("/judge/login")]
pub fn judge_login() -> Template {
    Template::render("judge-login", context! { title: get_title() })
}

#[rocket::get("/judge")]
pub fn judge() -> Template {
    Template::render("judge", context! { title: get_title() })
}

#[rocket::get("/judge/welcome")]
pub fn judge_welcome() -> Template {
    Template::render("judge-welcome", context! { title: get_title() })
}

#[rocket::get("/admin/login")]
pub fn admin_login() -> Template {
    Template::render("admin-login", context! { title: get_title() })
}

#[rocket::get("/admin")]
pub fn admin() -> Template {
    Template::render("admin", context! { title: get_title() })
}

fn get_title() -> String {
    env::var("GAVEL_NAME").expect("GAVEL_NAME not defined")
}
