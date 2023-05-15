use bson::doc;
use mongodb::error::Error;
use mongodb::Database;
use futures::stream::{TryStreamExt};

use super::models::{Judge, Project};
use crate::util::types::Stats;

pub async fn aggregate_stats(db: &Database) -> Result<Stats, Error> {
    // Fetch the collections we will use
    let judges_col = db.collection::<Judge>("judges");
    let projects_col = db.collection::<Project>("projects");

    // Sum to calculate the # of judges + projects
    let judges = judges_col.estimated_document_count(None).await?;
    let projects = projects_col.estimated_document_count(None).await?;

    // Aggregation for total seen
    let seen_pipeline = vec![doc! {
        "$group": {
            "_id": null,
            "total_seen": {
                "$sum": "$seen"
            }
        }
    }];
    let mut seen_cursor = projects_col.aggregate(seen_pipeline, None).await?;
    let mut seen = 0;
    if seen_cursor.advance().await? {
        seen = seen_cursor
            .current()
            .get_i32("total_seen")
            .unwrap_or_else(|_| 0);
    }

    // Aggregation for total votes
    let seen_pipeline = vec![doc! {
        "$group": {
            "_id": null,
            "total_votes": {
                "$sum": "$votes"
            }
        }
    }];
    let mut votes_cursor = projects_col.aggregate(seen_pipeline, None).await?;
    let mut votes = 0;
    if votes_cursor.advance().await? {
        votes = votes_cursor
            .current()
            .get_i32("total_seen")
            .unwrap_or_else(|_| 0);
    }

    // Aggregation for average mu
    let seen_pipeline = vec![doc! {
        "$group": {
            "_id": null,
            "avg_mu": {
                "$avg": "$mu"
            }
        }
    }];
    let mut mu_cursor = projects_col.aggregate(seen_pipeline, None).await?;
    let mut avg_mu = 0.0;
    if mu_cursor.advance().await? {
        avg_mu = mu_cursor
            .current()
            .get_f64("avg_mu")
            .unwrap_or_else(|_| 0.0);
    }

    // Aggregation for average sigma squared
    let seen_pipeline = vec![doc! {
        "$group": {
            "_id": null,
            "avg_sigma": {
                "$avg": "$sigma_sq"
            }
        }
    }];
    let mut mu_cursor = projects_col.aggregate(seen_pipeline, None).await?;
    let mut avg_sigma = 0.0;
    if mu_cursor.advance().await? {
        avg_sigma = mu_cursor
            .current()
            .get_f64("avg_sigma")
            .unwrap_or_else(|_| 0.0);
    }

    Ok(Stats {
        projects,
        seen: seen.try_into().unwrap_or_else(|_| 0),
        votes: votes.try_into().unwrap_or_else(|_| 0),
        time: 0, // TODO: find a way to store this in rocket's managed state
        avg_mu,
        avg_sigma,
        judges,
    })
}

pub async fn insert_projects(db: &Database, projects: Vec<Project>) -> Result<(), Error> {
    let collection = db.collection::<Project>("projects");
    collection.insert_many(projects, None).await?;
    Ok(())
}

pub async fn find_all_projects(db: &Database) -> Result<Vec<Project>, Error> {
    let collection = db.collection::<Project>("projects");
    let cursor = collection.find(None, None).await?;
    let projects = cursor.try_collect().await?;
    Ok(projects)
}
