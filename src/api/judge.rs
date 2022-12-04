use bson::doc;
use chrono::prelude::*;
use mongodb::{Collection, Database};
use rand::{distributions::Alphanumeric, Rng};
use rocket::http::Status;
use rocket::serde::{json::Json, Deserialize};
use rocket::State;

use crate::db::models::Judge;
use crate::util::crowd_bt;

#[derive(Deserialize)]
#[serde()]
pub struct Login<'r> {
    code: &'r str,
}

#[derive(Deserialize)]
#[serde()]
pub struct NewJudge<'r> {
    name: &'r str,
    email: &'r str,
    notes: &'r str,
}

#[rocket::post("/judge/login", data = "<body>")]
pub async fn login(db: &State<Database>, body: Json<Login<'_>>) -> (Status, String) {
    // Check in DB for correct code
    let collection: Collection<Judge> = db.collection("judges");
    let doc: Result<Option<Judge>, _> = collection
        .find_one(Some(doc! { "code": body.code }), None)
        .await;

    // If code is invalid or DB access fails, return error
    let judge = match doc {
        Ok(content) => match content {
            Some(x) => x,
            None => {
                return (Status::BadRequest, "Invalid code".to_string());
            }
        },
        Err(_) => {
            return (
                Status::InternalServerError,
                "Unable to fetch from database".to_string(),
            )
        }
    };

    // Generate random token
    let token: String = rand::thread_rng()
        .sample_iter(&Alphanumeric)
        .take(16)
        .map(char::from)
        .collect();

    // Update the token in the database
    let collection: Collection<Judge> = db.collection("judges");
    match collection
        .update_one(
            doc! {"code": judge.code},
            doc! {"$set": {"token": token.clone() }},
            None,
        )
        .await
    {
        Ok(_) => (Status::Ok, (format!("{}", token))),
        Err(_) => (
            Status::InternalServerError,
            "Unable to create token".to_string(),
        ),
    }
}

#[rocket::post("/judge/new", data = "<body>")]
pub async fn new_judge(db: &State<Database>, body: Json<NewJudge<'_>>) -> (Status, String) {
    // Create the new judge
    let judge = Judge {
        id: None,
        code: rand::thread_rng().gen_range(100000..999999).to_string(),
        token: "".to_string(),
        name: body.name.to_string(),
        email: body.email.to_string(),
        active: true,
        last_activity: Utc::now(),
        read_welcome: false,
        notes: body.notes.to_string(),
        next: None,
        prev: None,
        alpha: crowd_bt::ALPHA_PRIOR,
        beta: crowd_bt::BETA_PRIOR,
    };

    // Insert into database
    let collection: Collection<Judge> = db.collection("judges");
    let result: Result<_, _> = collection.insert_one(judge, None).await;
    match result {
        Ok(_) => (Status::Accepted, "{}".to_string()),
        Err(_) => (Status::BadRequest, "Invalid code".to_string()),
    }
}
