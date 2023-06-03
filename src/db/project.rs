use bson::{doc, oid::ObjectId};
use futures::stream::TryStreamExt;
use mongodb::{error::Error, Database};

use crate::api::request_types::NewProject;

use super::models::Project;

pub async fn insert_projects(db: &Database, projects: Vec<Project>) -> Result<(), Error> {
    let collection = db.collection::<Project>("projects");
    collection.insert_many(projects, None).await?;
    Ok(())
}

pub async fn insert_project(db: &Database, project: NewProject) -> Result<(), Error> {
    let project: Project = project.into();
    let collection = db.collection::<Project>("projects");
    collection.insert_one(project, None).await?;
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
    println!("Deleted {} documents", res.deleted_count);
    if res.deleted_count == 0 {
        return Err("No documents deleted".into());
    }
    Ok(())
}
