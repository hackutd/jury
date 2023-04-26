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

pub fn get_cookie<'r>(request: &'r rocket::Request<'_>, name: &str) -> Option<&'r Cookie<'static>> {
    request.cookies().get(name)
}

pub async fn parse_csv_from_data(csv: Data<'_>) -> String {
    // Extract all data and store into a string
    // TODO: No unwrap
    let str = csv.open(16.mebibytes()).into_string().await.unwrap();

    // Find the start and end of the CSV data
    // and cut out the CSV data from the raw string
    let start = str.find(CSV_HEADER).unwrap() + CSV_HEADER.len();
    let end = str[start..].find(CSV_FOOTER).unwrap();
    let cut_str = &str[start..(start + end)];

    return cut_str.to_string();
}
