use std::error::Error;

use mongodb::Database;

use super::models::Project;

pub async fn insert_projects(db: &Database, projects: Vec<Project>) -> Result<(), Box<dyn Error>> {
    let collection = db.collection::<Project>("projects");
    collection.insert_many(projects, None).await?;
    Ok(())
}
