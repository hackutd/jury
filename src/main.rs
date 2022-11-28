// https://docs.rs/dotenv/latest/dotenv/
use dotenv::dotenv;
use tokio;

use gavel3::db;

#[tokio::main]
async fn main() {
    // Load the .env file
    dotenv().ok();

    // Initialize database
    let init_result = db::init::init_db();
    match init_result.await {
        Ok(_) => {println!("Connected to MongoDB successfully!")},
        Err(e) => {println!("Error connecting to MongoDB {e}")},
    }
}
