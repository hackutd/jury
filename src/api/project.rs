use mongodb::Database;
use rocket::form::Form;
use rocket::http::Status;
use rocket::serde::json::Json;
use rocket::State;
use serde::Serialize;
use std::sync::Arc;

use super::request_types::NewProject;
use super::util::{AdminPassword, Token};
use crate::db::models::Project;
use crate::db::project::{
    aggregate_project_stats, delete_project_by_id, find_all_projects, find_project_by_id_string,
    insert_project,
};
use crate::try_status;
use crate::util::parse_csv::devpost_integration;
use crate::util::types::{CsvUpload, ProjectStats};
use crate::{db::project::insert_projects, util::parse_csv::parse_projects_csv};

#[rocket::post("/project/devpost", data = "<upload>")]
pub async fn add_devpost_csv(
    upload: Form<CsvUpload>,
    db: &State<Arc<Database>>,
    _password: AdminPassword,
) -> (Status, String) {
    // Parse the CSV data
    let parsed_csv = try_status!(
        devpost_integration(upload.csv.clone(), db).await,
        "Unable to parse CSV",
        Status::BadRequest
    );

    // Save the parsed CSV data to the database
    try_status!(
        insert_projects(&db, parsed_csv).await,
        "Unable to insert projects into database",
        Status::InternalServerError
    );

    (Status::Ok, "".to_string())
}

#[rocket::post("/project/new", data = "<body>")]
pub async fn new_project(
    db: &State<Arc<Database>>,
    body: Json<NewProject>,
    _password: AdminPassword,
) -> (Status, String) {
    match insert_project(db, body.0).await {
        Ok(_) => (Status::Accepted, "{}".to_string()),
        Err(e) => (
            Status::InternalServerError,
            format!("Unable to insert project: {}", e),
        ),
    }
}

#[derive(Serialize)]
pub struct ProjectCsvResponse {
    pub name: String,
    pub description: String,
    pub try_link: Option<String>,
    pub video_link: Option<String>,
    pub challenge_list: Vec<String>,
}

#[rocket::post("/project/csv", data = "<upload>")]
pub async fn preview_projects_csv(
    upload: Form<CsvUpload>,
    _password: AdminPassword,
) -> (Status, Json<Vec<ProjectCsvResponse>>) {
    // Parse the CSV data
    let projects = match parse_projects_csv(upload.csv.clone(), upload.has_header).await {
        Ok(x) => x,
        Err(_) => return (Status::BadRequest, Json(vec![])),
    };

    // Map projects to ProjectCsvResponse and return
    (
        Status::Ok,
        Json(
            projects
                .into_iter()
                .map(|x| ProjectCsvResponse {
                    name: x.name,
                    description: x.description,
                    try_link: x.try_link,
                    video_link: x.video_link,
                    challenge_list: x.challenge_list,
                })
                .collect::<Vec<ProjectCsvResponse>>(),
        ),
    )
}

#[rocket::post("/project/csv/upload", data = "<upload>")]
pub async fn add_projects_csv(
    upload: Form<CsvUpload>,
    db: &State<Arc<Database>>,
    _password: AdminPassword,
) -> (Status, String) {
    // Parse the CSV data
    let projects = try_status!(
        parse_projects_csv(upload.csv.clone(), upload.has_header).await,
        "Unable to parse CSV",
        Status::BadRequest
    );

    // If there are no projects, return Ok
    let num_projects = projects.len();
    if num_projects == 0 {
        return (Status::Ok, "0".to_string());
    }

    // Save the parsed CSV data to the database
    try_status!(
        insert_projects(db, projects).await,
        "Unable to insert projects into database",
        Status::InternalServerError
    );

    (Status::Ok, format!("{}", num_projects).to_string())
}

#[rocket::get("/project/list")]
pub async fn get_projects(
    db: &State<Arc<Database>>,
    _password: AdminPassword,
) -> (Status, Json<Vec<Project>>) {
    let project_list = match find_all_projects(db).await {
        Ok(p) => p,
        Err(e) => {
            eprintln!("Unable to get all projects: {e}",);
            return (Status::InternalServerError, Json(Vec::new()));
        }
    };

    (Status::Ok, Json(project_list))
}

#[rocket::delete("/project/<id>")]
pub async fn delete_project(
    db: &State<Arc<Database>>,
    id: &str,
    _password: AdminPassword,
) -> (Status, String) {
    match delete_project_by_id(db, id).await {
        Ok(_) => (Status::Ok, "{}".to_string()),
        Err(e) => (
            Status::InternalServerError,
            format!("Unable to delete project: {}", e),
        ),
    }
}

#[rocket::get("/project/stats")]
pub async fn project_stats(
    db: &State<Arc<Database>>,
    _password: AdminPassword,
) -> (Status, Json<ProjectStats>) {
    match aggregate_project_stats(db).await {
        Ok(stats) => (Status::Ok, Json(stats)),
        Err(_) => (Status::InternalServerError, Json(ProjectStats::default())),
    }
}

#[rocket::get("/project/<id>")]
pub async fn get_project_by_id(
    db: &State<Arc<Database>>,
    id: &str,
    _token: Token,
) -> (Status, Json<Project>) {
    match find_project_by_id_string(db, id).await {
        Ok(p) => (Status::Ok, Json(p)),
        Err(_) => (Status::InternalServerError, Json(Project::default())),
    }
}
