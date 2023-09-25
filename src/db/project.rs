use bson::{doc, oid::ObjectId};
use chrono::Utc;
use futures::stream::TryStreamExt;
use mongodb::{error::Error, Database};

use crate::{api::request_types::NewProject, util::types::ProjectStats};

use super::{
    models::{Judge, Project},
    options::{get_next_location, increment_location},
};

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

pub async fn find_all_active_projects(db: &Database) -> Result<Vec<Project>, Error> {
    let collection = db.collection::<Project>("projects");
    let cursor = collection.find(doc! { "active": true }, None).await?;
    let projects = cursor.try_collect().await?;
    Ok(projects)
}

/// Find all projects that are busy.
/// A project is considered busy if it has a next judge
pub async fn find_all_busy_projects(db: &Database) -> Result<Vec<ObjectId>, Error> {
    // Get list of busy project IDs
    let judge_collection = db.collection::<Judge>("judges");
    let cursor = judge_collection.find(
        doc! { "active": true, "next": { "$ne": None::<ObjectId> } },
        None,
    );
    let judges = cursor.await?.try_collect::<Vec<Judge>>().await?;
    let busy_projects = judges
        .into_iter()
        .map(|j| j.next.unwrap_or(ObjectId::default()))
        .collect::<Vec<ObjectId>>();

    Ok(busy_projects)
}

pub async fn find_project_by_id_string(
    db: &Database,
    id: &str,
) -> Result<Project, Box<dyn std::error::Error>> {
    let id = ObjectId::parse_str(id)?;
    find_project_by_id(db, id).await
}

/// Finds a project by its ID
pub async fn find_project_by_id(
    db: &Database,
    id: ObjectId,
) -> Result<Project, Box<dyn std::error::Error>> {
    let collection = db.collection::<Project>("projects");
    let project = collection.find_one(doc! { "_id": id }, None).await?;
    match project {
        Some(project) => Ok(project),
        None => Err("No project found".into()),
    }
}

/// Finds multiple projects given a vector of IDs
pub async fn find_projects_by_id(
    db: &Database,
    ids: Vec<ObjectId>,
) -> Result<Vec<Project>, Box<dyn std::error::Error>> {
    let collection = db.collection::<Project>("projects");
    let cursor = collection.find(doc! { "_id": { "$in": ids }}, None).await?;
    let projects = cursor.try_collect().await?;
    Ok(projects)
}

/// Updates a project's mu and sigma^2 values
pub async fn update_project_after_vote(
    db: &Database,
    id: Option<ObjectId>,
    new_mu: f64,
    new_sigma_sq: f64,
    won: bool,
) -> Result<(), Box<dyn std::error::Error>> {
    let collection = db.collection::<Project>("projects");

    // Get project
    let project = collection
        .find_one(doc! { "_id": id }, None)
        .await?
        .ok_or_else(|| "Could not find project by ID when updating project")?;

    // Get current time
    let now = Utc::now();

    // Add 1 to votes if won
    let votes = project.votes + if won { 1 } else { 0 };

    // Update project
    collection
        .update_one(
            doc! { "_id": id },
            doc! { "$set": { "mu": new_mu, "sigma_sq": new_sigma_sq, "votes": votes, "seen": project.seen + 1, "last_activity": now } },
            None,
        )
        .await?;

    Ok(())
}
