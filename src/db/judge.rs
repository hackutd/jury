use bson::{doc, oid::ObjectId};
use futures::stream::TryStreamExt;
use mongodb::{error::Error, Collection, Database};
use rocket::http::Status;

use crate::{api::request_types::NewJudge, util::types::JudgeStats};

use super::models::Judge;

pub async fn insert_judge(db: &Database, judge: NewJudge) -> Result<(), Error> {
    // Create the judge
    let judge: Judge = judge.into();

    // Insert into database
    let collection: Collection<Judge> = db.collection("judges");
    collection.insert_one(judge, None).await?;
    Ok(())
}

pub async fn insert_judges(db: &Database, judges: Vec<Judge>) -> Result<(), Error> {
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

pub async fn update_judge_token(db: &Database, code: &str, token: &str) -> Result<(), Error> {
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

pub async fn read_welcome(db: &Database, token: &str) -> Result<(), Error> {
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

pub async fn delete_judge_by_id(
    db: &Database,
    id: &str,
) -> Result<(), Box<dyn std::error::Error>> {
    let collection = db.collection::<Judge>("judges");
    let res = collection
        .delete_one(doc! { "_id": ObjectId::parse_str(id)? }, None)
        .await?;
    if res.deleted_count == 0 {
        return Err("No documents deleted".into());
    }
    Ok(())
}

pub async fn aggregate_judge_stats(db: &Database) -> Result<JudgeStats, mongodb::error::Error> {
    let judges_col = db.collection::<Judge>("judges");
    let num = judges_col.estimated_document_count(None).await?;

    // Aggregation for count of active judges and the average views among them
    let pipeline = vec![
        doc! {
            "$match": {
                "active": true
            }
        },
        doc! {
            "$group": {
                "_id": null,
                "count": {
                    "$sum": 1
                },
                "votes": {
                    "$avg": "$votes"
                }
            }
        },
    ];
    let mut cursor = judges_col.aggregate(pipeline, None).await?;
    let mut count = 0;
    let mut votes = 0.0;
    if cursor.advance().await? {
        count = cursor.current().get_i32("count").unwrap_or_else(|_| 0);
        // votes = cursor.current().get_f64("votes").unwrap_or_else(|_| 0.0);
        votes = cursor.current().get_f64("votes").unwrap();
    }

    Ok(JudgeStats {
        num,
        avg_votes: votes,
        num_active: count.try_into().unwrap_or(0),
    })
}

pub async fn find_all_judges(db: &Database) -> Result<Vec<Judge>, Error> {
    let collection = db.collection::<Judge>("judges");
    let cursor = collection.find(None, None).await?;
    let judges = cursor.try_collect().await?;
    Ok(judges)
}
