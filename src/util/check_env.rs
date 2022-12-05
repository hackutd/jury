use std::env;

pub fn check() {
    let _ = env::var("MONGODB_URI").expect("MONGODB_URI not defined");
    let _ = env::var("GAVEL_NAME").expect("GAVEL_NAME not defined");
    let _ = env::var("GAVEL_ADMIN_PASSWORD").expect("GAVEL_ADMIN_PASSWORD not defined");
}