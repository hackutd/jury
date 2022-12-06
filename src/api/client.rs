use mongodb::Database;
use rocket::{http::CookieJar, response::Redirect, State};
use rocket_dyn_templates::{context, Template};

use std::env;

use crate::db::judge::find_judge_by_token;

use super::util::{AdminPassword, Token};

#[rocket::get("/")]
pub fn home() -> Template {
    Template::render("index", context! { title: get_title() })
}

// Will redirect to judge page if already logged in
#[rocket::get("/judge/login")]
pub async fn judge_login(
    db: &State<Database>,
    cookies: &CookieJar<'_>,
) -> Result<Redirect, Template> {
    // Create login template
    let template = Template::render("judge-login", context! { title: get_title() });

    // Get token from cookies
    let token = match cookies.get("token") {
        Some(t) => t,
        None => return Err(template),
    };

    // Find judge and return success or error
    match find_judge_by_token(db, token.value()).await {
        Ok(_) => Ok(Redirect::to("/judge")),
        Err(_) => Err(template),
    }
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
pub fn admin_login(cookies: &CookieJar<'_>) -> Result<Redirect, Template> {
    // Create login template
    let template = Template::render("admin-login", context! { title: get_title() });

    // Get token from cookies
    let pass = match cookies.get("admin-pass") {
        Some(t) => t,
        None => return Err(template),
    };

    // Find judge and return success or error
    if pass.value() == env::var("GAVEL_ADMIN_PASSWORD").expect("GAVEL_ADMIN_PASSWORD not defined") {
        Ok(Redirect::to("/admin"))
    } else {
        Err(template)
    }
}

#[rocket::get("/admin")]
pub fn admin(_password: AdminPassword) -> Template {
    Template::render("admin", context! { title: get_title() })
}

fn get_title() -> String {
    env::var("GAVEL_NAME").expect("GAVEL_NAME not defined")
}
