use bson::{doc, oid::ObjectId};
use futures::stream::TryStreamExt;
use mongodb::{error::Error, Database};

use crate::{api::request_types::NewProject, util::types::ProjectStats};

use super::{models::Project, options::{get_next_location, increment_location}};

pub async fn insert_projects(db: &Database, mut projects: Vec<Project>) -> Result<(), Error> {
    let new_location = get_next_location(db).await?;
    let collection = db.collection::<Project>("projects");

    // Update the location of all projects
    for (i, project) in projects.iter_mut().enumerate() {
        project.location = new_location + i as u64;
    }

    let proj_len = projects.len() as u32;
    collection.insert_many(projects, None).await?;
    increment_location(db, Some(proj_len)).await?;
    Ok(())
}

pub async fn insert_project(db: &Database, project: NewProject) -> Result<(), Error> {
    let new_location = get_next_location(db).await?;
    let mut project: Project = project.into();
    project.location = new_location;
    let collection = db.collection::<Project>("projects");
    // TODO: This should be a transaction
    collection.insert_one(project, None).await?;
    increment_location(db, None).await?;
    Ok(())
}

pub async fn find_all_projects(db: &Database) -> Result<Vec<Project>, Error> {
    let collection = db.collection::<Project>("projects");
    let cursor = collection.find(None, None).await?;
    let projects = cursor.try_collect().await?;
    Ok(projects)
}

pub async fn delete_project_by_id(
    db: &Database,
    id: &str,
) -> Result<(), Box<dyn std::error::Error>> {
    let collection = db.collection::<Project>("projects");
    let res = collection
        .delete_one(doc! { "_id": ObjectId::parse_str(id)? }, None)
        .await?;
    if res.deleted_count == 0 {
        return Err("No documents deleted".into());
    }
    Ok(())
}

pub async fn aggregate_project_stats(db: &Database) -> Result<ProjectStats, mongodb::error::Error> {
    let projects_col = db.collection::<Project>("projects");

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
                },
                "seen": {
                    "$avg": "$seen"
                }
            }
        },
    ];
    let mut cursor = projects_col.aggregate(pipeline, None).await?;
    let mut count = 0;
    let mut votes = 0.0;
    let mut seen = 0.0;
    if cursor.advance().await? {
        count = cursor.current().get_i32("count").unwrap_or_else(|_| 0);
        // votes = cursor.current().get_f64("votes").unwrap_or_else(|_| 0.0);
        votes = cursor.current().get_f64("votes").unwrap();
        seen = cursor.current().get_f64("seen").unwrap();
    }

    Ok(ProjectStats {
        num: count.try_into().unwrap_or(0),
        avg_votes: votes,
        avg_seen: seen,
    })
}
