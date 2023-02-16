use std::env;

pub fn check() {
    let _ = env::var("MONGODB_URI").expect("MONGODB_URI not defined");
    let _ = env::var("JURY_NAME").expect("JURY_NAME not defined");
    let _ = env::var("JURY_URL").expect("JURY_URL not defined");
    let _ = env::var("JURY_ADMIN_PASSWORD").expect("JURY_ADMIN_PASSWORD not defined");
}