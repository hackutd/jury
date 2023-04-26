use std::error::Error;

use bson::doc;
use mongodb::{Collection, Database};
use rocket::http::Status;

use crate::api::request_types::NewJudge;

use super::models::Judge;

pub async fn insert_judge(db: &Database, body: NewJudge<'_>) -> Result<(), Status> {
    // Create the new judge
    let judge = Judge::new(
        body.name.to_string(),
        body.email.to_string(),
        body.notes.to_string(),
    );

    // Insert into database
    let collection: Collection<Judge> = db.collection("judges");
    let result: Result<_, _> = collection.insert_one(judge, None).await;
    match result {
        Ok(_) => Ok(()),
        Err(_) => Err(Status::InternalServerError),
    }
}

pub async fn insert_judges(db: &Database, judges: Vec<Judge>) -> Result<(), Box<dyn Error>> {
    let collection = db.collection::<Judge>("judges");
    collection.insert_many(judges, None).await?;
    Ok(())
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
    let doc: Result<Option<Judge>, _> = collection
        .find_one(Some(doc! { "token": token }), None)
        .await;

    // If code is invalid or DB access fails, return error
    match doc {
        Ok(content) => match content {
            Some(x) => Ok(x),
            None => Err(Status::Unauthorized),
        },
        Err(_) => Err(Status::InternalServerError),
    }
}

pub async fn update_judge_token(
    db: &Database,
    code: &str,
    token: &str,
) -> Result<(), Box<dyn Error>> {
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

pub async fn read_welcome(db: &Database, token: &str) -> Result<(), Box<dyn Error>> {
    let collection: Collection<Judge> = db.collection("judges");
    collection
        .update_one(
            doc! {"token": token},
            doc! {"$set": {"read_welcome": true}},
            None,
        )
        .await?;

    Ok(())
}
