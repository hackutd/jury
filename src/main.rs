use dotenv::dotenv;
use rocket::fs::{relative, FileServer};
use rocket_dyn_templates::Template;

use jury::api::{client, judge, admin, catchers};
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
                client::judge_login,
                client::judge,
                client::judge_welcome,
                client::admin_login,
                client::admin,
                judge::login,
                judge::new_judge,
                judge::judge_read_welcome,
                admin::login
            ],
        )
        .register("/", catchers![catchers::unauthorized])
        .mount("/static", FileServer::from(relative!("public")))
        .attach(Template::fairing())
}
