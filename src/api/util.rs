use std::error::Error;

use rocket::data::{Data, ToByteUnit};
use rocket::http::Cookie;

const CSV_HEADER: &str = "Content-Type: text/csv\r\n\r\n";
const CSV_FOOTER: &str = "\r\n----------------------------";

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

pub async fn parse_csv_from_data(csv: Data<'_>) -> Result<String, Box<dyn Error>> {
    // Extract all data and store into a string
    let str = csv.open(16.mebibytes()).into_string().await?;

    println!("str: {}", str.to_string());

    // Find the start and end of the CSV data
    // and cut out the CSV data from the raw string
    let start = str
        .find(CSV_HEADER)
        .ok_or("Unable to find CSV header with correct content-type")?
        + CSV_HEADER.len();
    let end = str[start..]
        .find(CSV_FOOTER)
        .ok_or("Unable to find CSV footer")?;
    let cut_str = &str[start..(start + end)];

    return Ok(cut_str.to_string());
}
