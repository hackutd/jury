use rocket_dyn_templates::Template;

use dotenv::dotenv;
use gavel3::api::{client, server};
use gavel3::{db, util};

#[macro_use]
extern crate rocket;

#[launch]
async fn rocket() -> _ {
    // Load environmental variables from .env file
    dotenv().ok();

    // Check to make sure all env vars are there
    util::check_env::check();

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
    rocket::build()
        .mount(
            "/",
            routes![
                client::home,
                client::judge_login,
                client::judge,
                client::judge_welcome,
                client::admin_login,
                client::admin,
                server::judge_login
            ],
        )
        .attach(Template::fairing())
}
