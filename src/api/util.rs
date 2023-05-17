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

/// Macro to convert String to Option<String>.
/// If the string is empty, return None.
/// Otherwise, return Some(string).
#[macro_export]
macro_rules! str_opt {
    ($x:expr) => {
        if $x.is_empty() {
            None
        } else {
            Some($x.to_string())
        }
    };
}

/// Macro that parses a result returns a 500 error if it fails.
/// Otherwise, return the result of the function.
#[macro_export]
macro_rules! try_status {
    ($x:expr, $msg:literal, $status:expr) => {
        match $x {
            Ok(x) => x,
            Err(e) => {
                eprintln!("{}: {}", $msg, e);
                return ($status, format!("{}: {}", $msg, e));
            }
        }
    };
}

pub fn get_cookie<'r>(request: &'r rocket::Request<'_>, name: &str) -> Option<&'r Cookie<'static>> {
    request.cookies().get(name)
}
