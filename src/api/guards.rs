use async_trait::async_trait;
use mongodb::Database;
use rocket::http::Status;
use rocket::outcome::Outcome;
use rocket::request::FromRequest;
use std::env;

use crate::db::judge::find_judge_by_token;

use super::util::{get_cookie, AdminPassword, AdminPasswordError, Token, TokenError};

// Route guard for judge
#[async_trait]
impl<'r> FromRequest<'r> for Token {
    type Error = TokenError;
    async fn from_request(
        request: &'r rocket::Request<'_>,
    ) -> rocket::request::Outcome<Self, Self::Error> {
        // Retrieve token cookie
        let cookie_token = match get_cookie(request, "token") {
            Some(c) => c.value(),
            None => {
                return Outcome::Failure((Status::Unauthorized, TokenError::Missing));
            }
        };

        // Get database from saved state
        let db = match request.rocket().state::<Database>() {
            Some(d) => d,
            None => {
                return Outcome::Failure((
                    Status::InternalServerError,
                    TokenError::InternalServerError,
                ));
            }
        };

        // Find judge and return success or error
        match find_judge_by_token(db, cookie_token).await {
            Ok(_) => Outcome::Success(Token(cookie_token.to_string())),
            Err(status) => Outcome::Failure((status, TokenError::InternalServerError)),
        }
    }
}

// Route guard for admin
#[async_trait]
impl<'r> FromRequest<'r> for AdminPassword {
    type Error = AdminPasswordError;
    async fn from_request(
        request: &'r rocket::Request<'_>,
    ) -> rocket::request::Outcome<Self, Self::Error> {
        // Pull password from cookies
        let cookie_password = match request.cookies().get("admin-pass") {
            Some(c) => c.value(),
            None => {
                return Outcome::Failure((Status::Unauthorized, AdminPasswordError::Missing));
            }
        };

        // Get correct passsword from environmental variables
        let correct = env::var("JURY_ADMIN_PASSWORD").expect("JURY_ADMIN_PASSWORD not defined");

        // Compare
        if cookie_password == correct {
            Outcome::Success(AdminPassword(cookie_password.to_string()))
        } else {
            Outcome::Failure((Status::Unauthorized, AdminPasswordError::Invalid))
        }
    }
}
