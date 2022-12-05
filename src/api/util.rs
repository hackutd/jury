use rocket::http::Cookie;

pub struct Token(pub String);

#[derive(Debug)]
pub enum TokenError {
    Invalid,
    Missing,
    InternalServerError,
}

pub struct AdminPassword(pub String);

#[derive(Debug)]
pub enum AdminPasswordError {
    Invalid,
    Missing,
}

pub fn get_cookie<'r>(request: &'r rocket::Request<'_>, name: &str) -> Option<&'r Cookie<'static>> {
    request.cookies().get(name)
}