use dotenv::dotenv;
use jury::api::client::CORS;
use jury::util::tasks::init_threads;
use rocket::fs::{relative, FileServer};
use std::env;
use std::sync::Arc;

use jury::api::{admin, client, judge, project};
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
    let db = Arc::new(init_result.await.expect("Could not connect to MongoDB!"));

    // Select FileServer
    let fileserver_path = env::var("FILESERVER").unwrap_or_else(|_| "".to_string());
    let files = if fileserver_path.starts_with("/") {
        FileServer::from(fileserver_path)
    } else {
        FileServer::from(relative!("public"))
    };

    // Initialize all threads
    // and get sender for sending new message senders
    // See https://cdn.michaelzhao.xyz/archive/jury-network.png
    let vector_tx = init_threads(db.clone()).await;

    // Start server
    rocket::build()
        .manage(db.clone())
        .manage(vector_tx)
        .mount(
            "/api",
            routes![
                judge::login,
                judge::new_judge,
                judge::judge_read_welcome,
                judge::add_judges_csv,
                judge::judge_stats,
                project::add_devpost_csv,
                admin::login,
                admin::get_stats,
                admin::req_sync,
            ],
        )
        .mount("/", routes![client::home, client::all_options])
        .mount("/", files)
        .attach(CORS)
}
