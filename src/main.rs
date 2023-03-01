use dotenv::dotenv;
use jury::api::client::CORS;
use rocket::fs::{relative, FileServer};
use std::env;

use jury::api::{admin, client, judge};
use jury::{db, util};

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
    let db = init_result.await.expect("Could not connect to MongoDB!");

    // Select FileServer
    let fileserver_path = env::var("FILESERVER").unwrap_or_else(|_| "".to_string());
    let files = if fileserver_path.starts_with("/") {
        FileServer::from(fileserver_path)
    } else {
        FileServer::from(relative!("public"))
    };

    println!("ur mum");

    // Start server
    rocket::build()
        .manage(db)
        .mount(
            "/api",
            routes![
                judge::login,
                judge::new_judge,
                judge::judge_read_welcome,
                admin::login,
                admin::get_stats,
            ],
        )
        .mount("/", routes![client::home, client::all_options])
        .mount("/", files)
        .attach(CORS)
}
