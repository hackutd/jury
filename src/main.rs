use dotenv::dotenv;
use jury::api::client::CORS;
use rocket::fs::{relative, FileServer};

use jury::api::{admin, catchers, client, judge};
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

    // Start server
    rocket::build()
        .manage(db)
        .mount(
            "/",
            routes![
                client::home,
                judge::login,
                judge::new_judge,
                judge::judge_read_welcome,
                admin::login,
                client::all_options
            ],
        )
        .register("/", catchers![catchers::unauthorized])
        .mount("/", FileServer::from(relative!("public")))
        .attach(CORS)
}
