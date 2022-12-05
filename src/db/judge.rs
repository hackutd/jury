use std::error::Error;

use bson::doc;
use chrono::Utc;
use mongodb::{Collection, Database};
use rand::Rng;
use rocket::http::Status;

use crate::{api::request_types::NewJudge, util::crowd_bt};

use super::models::Judge;

pub async fn insert_judge(db: &Database, body: NewJudge<'_>) -> Result<(), Status> {
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
        Ok(_) => Ok(()),
        Err(_) => Err(Status::InternalServerError),
    }
}

pub async fn find_judge_by_code(db: &Database, code: &str) -> Result<Judge, Status> {
    // Check in DB for correct code
    let collection: Collection<Judge> = db.collection("judges");
    let doc: Result<Option<Judge>, _> =
        collection.find_one(Some(doc! { "code": code }), None).await;

    // If code is invalid or DB access fails, return error
    match doc {
        Ok(content) => match content {
            Some(x) => Ok(x),
            None => Err(Status::BadRequest),
        },
        Err(_) => Err(Status::InternalServerError),
    }
}

pub async fn find_judge_by_token(db: &Database, token: &str) -> Result<Judge, Status> {
    // Check in DB for correct code
    let collection: Collection<Judge> = db.collection("judges");
    let doc: Result<Option<Judge>, _> =
        collection.find_one(Some(doc! { "token": token }), None).await;

    // If code is invalid or DB access fails, return error
    match doc {
        Ok(content) => match content {
            Some(x) => Ok(x),
            None => Err(Status::Unauthorized),
        },
        Err(_) => Err(Status::InternalServerError),
    }
}

pub async fn update_judge_token(db: &Database, code: &str, token: &str) -> Result<(), Box<dyn Error>> {
    let collection: Collection<Judge> = db.collection("judges");
    collection
        .update_one(
            doc! {"code": code},
            doc! {"$set": {"token": token.clone() }},
            None,
        )
        .await?;

    Ok(())
}
