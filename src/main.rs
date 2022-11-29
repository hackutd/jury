use dotenv::dotenv;

use gavel3::api;
use gavel3::db;

#[macro_use] extern crate rocket;

#[launch]
async fn rocket() -> _ {
    // Load the .env file
    dotenv().ok();

    // Initialize database
    let init_result = db::init::init_db();
    match init_result.await {
        Ok(_) => {
            println!("Connected to MongoDB successfully!")
        }
        Err(e) => {
            println!("Error connecting to MongoDB {e}")
        }
    }

    // Start server
    rocket::build().mount("/", routes![api::server::home])
}
