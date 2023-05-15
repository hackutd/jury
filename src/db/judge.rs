use std::error::Error;

use bson::doc;
use chrono::{Utc, TimeZone};
use mongodb::{Collection, Database};
use rocket::http::Status;

use crate::{api::request_types::NewJudge, util::types::JudgeStats};

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
        // TODO: This should be fine but if for some reason utc 0 is invalid, this should be changed lmao
        last_activity: Utc.timestamp_opt(0, 0).unwrap(),
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

pub async fn aggregate_judge_stats(db: &Database) -> Result<JudgeStats, mongodb::error::Error> {
    let judges_col = db.collection::<Judge>("judges");
    let num = judges_col.estimated_document_count(None).await?;

    // Aggregation for average alpha and beta
    let alpha_pipeline = vec![doc! {
        "$group": {
            "_id": null,
            "alpha": {
                "$avg": "$alpha"
            },
            "beta": {
                "$avg": "$beta"
            }
        }
    }];
    let mut cursor = judges_col.aggregate(alpha_pipeline, None).await?;
    let mut alpha = 0.0;
    let mut beta = 0.0;
    if cursor.advance().await? {
        alpha = cursor.current().get_f64("alpha").unwrap_or_else(|_| 0.0);
        beta = cursor.current().get_f64("beta").unwrap_or_else(|_| 0.0);
    }

    Ok(JudgeStats { num, alpha, beta })
}
