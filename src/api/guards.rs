use async_trait::async_trait;
use bson::doc;
use mongodb::Database;
use rocket::http::Status;
use rocket::outcome::Outcome;
use rocket::request::FromRequest;

use crate::db::models::Judge;

pub struct Token(String);

#[derive(Debug)]
pub enum TokenError {
    Invalid,
    Missing,
    InternalServerError,
}

// Route guard for judge
#[async_trait]
impl<'r> FromRequest<'r> for Token {
    type Error = TokenError;
    async fn from_request(
        request: &'r rocket::Request<'_>,
    ) -> rocket::request::Outcome<Self, Self::Error> {
        let cookie_token = match request.cookies().get("token") {
            Some(c) => c.value(),
            None => {
                return Outcome::Failure((Status::Unauthorized, TokenError::Missing));
            }
        };

        let db = match request.rocket().state::<Database>() {
            Some(d) => d,
            None => {
                return Outcome::Failure((
                    Status::InternalServerError,
                    TokenError::InternalServerError,
                ));
            }
        };

        let doc: Result<Option<Judge>, _> = db
            .collection("judges")
            .find_one(Some(doc! { "token": cookie_token }), None)
            .await;
        match doc {
            Ok(content) => match content {
                Some(_) => Outcome::Success(Token(cookie_token.to_string())),
                None => Outcome::Failure((Status::Unauthorized, TokenError::Invalid)),
            },
            Err(_) => {
                Outcome::Failure((Status::InternalServerError, TokenError::InternalServerError))
            }
        }
    }
}
