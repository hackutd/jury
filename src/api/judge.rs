use std::sync::Arc;

use bson::doc;
use bson::oid::ObjectId;
use mongodb::Database;
use rand::{distributions::Alphanumeric, Rng};
use rocket::form::Form;
use rocket::http::Status;
use rocket::serde::json::Json;
use rocket::State;
use serde::Serialize;

use crate::db::judge::{
    aggregate_judge_stats, delete_judge_by_id, find_all_judges, find_judge_by_token,
    update_judge_projects,
};
use crate::db::models::{Judge, Project};
use crate::db::project::{find_project_by_id, update_project_after_vote, find_projects_by_id};
use crate::util::crowd_bt::update;
use crate::util::judging_flow::pick_next_project;
use crate::util::types::{
    BooleanResponse, CsvUpload, JudgeNextProject, JudgeStats, JudgeVoteProjectInfo,
};
use crate::{
    db::judge::{
        find_judge_by_code, insert_judge, insert_judges, read_welcome, update_judge_token,
    },
    try_status,
    util::parse_csv::parse_judge_csv,
};

use super::request_types::JudgeVote;
use super::{
    request_types::{Login, NewJudge},
    util::{AdminPassword, Token},
};

#[rocket::get("/judge")]
pub async fn get_judge(db: &State<Arc<Database>>, token: Token) -> (Status, Json<Judge>) {
    let judge = match find_judge_by_token(db, &token.0).await {
        Ok(j) => j,
        Err(_) => return (Status::NotFound, Json(Judge::default())),
    };

    (Status::Ok, Json(judge))
}

#[rocket::post("/judge/login", data = "<body>")]
pub async fn login(db: &State<Arc<Database>>, body: Json<Login<'_>>) -> (Status, String) {
    // Find judge from db using code
    let judge = try_status!(
        find_judge_by_code(db, body.code).await,
        "Unable to process or find code",
        Status::BadRequest
    );

    // Generate random 16-character token
    let token: String = rand::thread_rng()
        .sample_iter(&Alphanumeric)
        .take(16)
        .map(char::from)
        .collect();

    // Update the token in the database and return token
    match update_judge_token(db, &judge.code, &token).await {
        Ok(_) => (Status::Ok, (format!("{}", token))),
        Err(_) => (
            Status::InternalServerError,
            "Unable to create token".to_string(),
        ),
    }
}

#[rocket::post("/judge/auth")]
pub async fn auth_judge(_token: Token) -> Status {
    // Request guard will fail if token is invalid/missing
    Status::Ok
}

#[rocket::post("/judge/new", data = "<body>")]
pub async fn new_judge(
    db: &State<Arc<Database>>,
    body: Json<NewJudge>,
    _password: AdminPassword,
) -> (Status, String) {
    match insert_judge(db, body.0).await {
        Ok(_) => (Status::Accepted, "{}".to_string()),
        Err(e) => (
            Status::InternalServerError,
            format!("Unable to insert judge: {}", e),
        ),
    }
}

#[derive(Serialize)]
pub struct JudgeCsvResponse {
    pub name: String,
    pub email: String,
    pub notes: String,
}

#[rocket::post("/judge/csv", data = "<upload>")]
pub async fn preview_judges_csv(
    upload: Form<CsvUpload>,
    _password: AdminPassword,
) -> (Status, Json<Vec<JudgeCsvResponse>>) {
    // Parse the CSV data
    let judges = match parse_judge_csv(upload.csv.clone(), upload.has_header).await {
        Ok(x) => x,
        Err(_) => return (Status::BadRequest, Json(vec![])),
    };

    // Map judges to JudgeCsvResponse and return
    (
        Status::Ok,
        Json(
            judges
                .into_iter()
                .map(|x| JudgeCsvResponse {
                    name: x.name,
                    email: x.email,
                    notes: x.notes,
                })
                .collect::<Vec<JudgeCsvResponse>>(),
        ),
    )
}

#[rocket::post("/judge/csv/upload", data = "<upload>")]
pub async fn add_judges_csv(
    upload: Form<CsvUpload>,
    db: &State<Arc<Database>>,
    _password: AdminPassword,
) -> (Status, String) {
    // Parse the CSV data
    let judges = try_status!(
        parse_judge_csv(upload.csv.clone(), upload.has_header).await,
        "Unable to parse CSV",
        Status::BadRequest
    );

    // If there are no judges, return Ok
    let num_judges = judges.len();
    if num_judges == 0 {
        return (Status::Ok, "0".to_string());
    }

    // Save the parsed CSV data to the database
    try_status!(
        insert_judges(db, judges).await,
        "Unable to insert judges into database",
        Status::InternalServerError
    );

    (Status::Ok, format!("{}", num_judges).to_string())
}

#[rocket::get("/judge/welcome")]
pub async fn check_judge_read_welcome(
    db: &State<Arc<Database>>,
    token: Token,
) -> (Status, Json<BooleanResponse>) {
    let judge = match find_judge_by_token(db, &token.0).await {
        Ok(j) => j,
        Err(e) => return (e, Json(BooleanResponse::new(false))),
    };

    (Status::Ok, Json(BooleanResponse::new(judge.read_welcome)))
}

#[rocket::post("/judge/welcome")]
pub async fn judge_read_welcome(db: &State<Arc<Database>>, token: Token) -> (Status, String) {
    match read_welcome(db, &token.0).await {
        Ok(_) => (Status::Accepted, "".to_string()),
        Err(_) => (
            Status::InternalServerError,
            "Internal Server Error".to_string(),
        ),
    }
}

#[rocket::get("/judge/stats")]
pub async fn judge_stats(
    db: &State<Arc<Database>>,
    _password: AdminPassword,
) -> (Status, Json<JudgeStats>) {
    match aggregate_judge_stats(db).await {
        Ok(stats) => (Status::Ok, Json(stats)),
        Err(_) => (Status::InternalServerError, Json(JudgeStats::default())),
    }
}

#[rocket::get("/judge/list")]
pub async fn get_judges(db: &State<Arc<Database>>) -> (Status, Json<Vec<Judge>>) {
    let judge_list = match find_all_judges(db).await {
        Ok(p) => p,
        Err(e) => {
            eprintln!("Unable to get all judges: {e}",);
            return (Status::InternalServerError, Json(Vec::new()));
        }
    };

    (Status::Ok, Json(judge_list))
}

#[rocket::delete("/judge/<id>")]
pub async fn delete_judge(
    db: &State<Arc<Database>>,
    id: &str,
    _password: AdminPassword,
) -> (Status, String) {
    match delete_judge_by_id(db, id).await {
        Ok(_) => (Status::Ok, "{}".to_string()),
        Err(e) => (
            Status::InternalServerError,
            format!("Unable to delete judge: {}", e),
        ),
    }
}

#[rocket::post("/judge/vote", data = "<body>")]
pub async fn judge_vote(
    db: &State<Arc<Database>>,
    token: Token,
    body: Json<JudgeVote>,
) -> (Status, Json<JudgeNextProject>) {
    // Get judge
    let judge = match find_judge_by_token(db, &token.0).await {
        Ok(j) => j,
        Err(e) => return (e, Json(JudgeNextProject::default())),
    };

    // If judge has no previous project, get the next and return ok
    if judge.prev.is_none() {
        // Find next project
        let next = match pick_next_project(db, &judge).await {
            Ok(p) => p,
            Err(e) => {
                eprintln!("{}", e);
                return (
                    Status::InternalServerError,
                    Json(JudgeNextProject::default()),
                );
            }
        };
        
        println!("NEXT PROJECT: {:?}", next.clone().unwrap_or(Project::default()).id);

        // Update judge's projects (update next/prev project) in database
        match update_judge_projects(db, &judge, &next, judge.alpha, judge.beta, false).await {
            Ok(()) => (),
            Err(e) => {
                eprintln!("{}", e);
                return (
                    Status::InternalServerError,
                    Json(JudgeNextProject::default()),
                );
            }
        };

        // Ok
        return (
            Status::Ok,
            Json(JudgeNextProject::new(
                judge.id.unwrap_or_else(|| ObjectId::new()).to_string(),
                judge.next.map(|p| p.to_string()),
                next.map(|p| {
                    p.id.map(|i| i.to_string())
                        .unwrap_or_else(|| "".to_string())
                }),
            )),
        );
    }

    // If judge DOES have previous project, we need to check their vote
    // and update everything accordingly

    // Get the two projects
    let prev_proj =
        match find_project_by_id(db, judge.prev.unwrap_or_else(|| ObjectId::new())).await {
            Ok(p) => p,
            Err(e) => {
                eprintln!("{}", e);
                return (
                    Status::InternalServerError,
                    Json(JudgeNextProject::default()),
                );
            }
        };
    let next_proj =
        match find_project_by_id(db, judge.next.unwrap_or_else(|| ObjectId::new())).await {
            Ok(p) => p,
            Err(e) => {
                eprintln!("{}", e);
                return (
                    Status::InternalServerError,
                    Json(JudgeNextProject::default()),
                );
            }
        };

    // Set winner/loser based on curr_winner
    let (winner_proj, loser_proj) = if body.0.curr_winner {
        (next_proj, prev_proj)
    } else {
        (prev_proj, next_proj)
    };

    // Run the update function
    let (n_alpha, n_beta, n_mu_winner, n_sigma_sq_winner, n_mu_loser, n_sigma_sq_loser) = update(
        judge.alpha,
        judge.beta,
        winner_proj.mu,
        winner_proj.sigma_sq,
        loser_proj.mu,
        loser_proj.sigma_sq,
    );

    // Update both projects
    match update_project_after_vote(db, winner_proj.id, n_mu_winner, n_sigma_sq_winner, true).await
    {
        Ok(_) => (),
        Err(e) => {
            eprintln!("{}", e);
            return (
                Status::InternalServerError,
                Json(JudgeNextProject::default()),
            );
        }
    }
    match update_project_after_vote(db, loser_proj.id, n_mu_loser, n_sigma_sq_loser, false).await {
        Ok(_) => (),
        Err(e) => {
            eprintln!("{}", e);
            return (
                Status::InternalServerError,
                Json(JudgeNextProject::default()),
            );
        }
    }

    // TODO: Create a table for storing votes

    // Create an updated judge
    let mut updated_judge = judge.clone();
    updated_judge.alpha = n_alpha;
    updated_judge.beta = n_beta;

    // Pick judge's next project
    let next = match pick_next_project(db, &judge).await {
        Ok(p) => p,
        Err(e) => {
            eprintln!("{}", e);
            return (
                Status::InternalServerError,
                Json(JudgeNextProject::default()),
            );
        }
    };

    // Update the judge
    match update_judge_projects(db, &judge, &next, n_alpha, n_beta, true).await {
        Ok(_) => (),
        Err(e) => {
            eprintln!("{}", e);
            return (
                Status::InternalServerError,
                Json(JudgeNextProject::default()),
            );
        }
    }

    println!(
        "{}, {}, {}",
        judge.id.unwrap_or_else(|| ObjectId::new()).to_string(),
        judge.next.map(|p| p.to_string()).unwrap_or("NULL".to_string()),
        next.clone().map(|p| {
            p.id.map(|i| i.to_string())
                .unwrap_or_else(|| "".to_string())
        }).unwrap_or("NULL".to_string()),
    );

    // Return
    (
        Status::Ok,
        Json(JudgeNextProject::new(
            judge.id.unwrap_or_else(|| ObjectId::new()).to_string(),
            judge.next.map(|p| p.to_string()),
            next.map(|p| {
                p.id.map(|i| i.to_string())
                    .unwrap_or_else(|| "".to_string())
            }),
        )),
    )
}

#[rocket::get("/judge/vote/info")]
pub async fn get_judge_voting_project_info(
    db: &State<Arc<Database>>,
    token: Token,
) -> (Status, Json<JudgeVoteProjectInfo>) {
    // Get judge by token
    let judge = match find_judge_by_token(db, &token.0).await {
        Ok(j) => j,
        Err(e) => return (e, Json(JudgeVoteProjectInfo::default())),
    };

    // Get the two projects
    let curr_proj = if judge.next.is_some() {
        match find_project_by_id(db, judge.next.unwrap()).await {
            Ok(p) => Some(p),
            Err(e) => {
                eprintln!("{}", e);
                return (
                    Status::InternalServerError,
                    Json(JudgeVoteProjectInfo::default()),
                );
            }
        }
    } else {
        None
    };
    let prev_proj = if judge.prev.is_some() {
        match find_project_by_id(db, judge.prev.unwrap()).await {
            Ok(p) => Some(p),
            Err(e) => {
                eprintln!("{}", e);
                return (
                    Status::InternalServerError,
                    Json(JudgeVoteProjectInfo::default()),
                );
            }
        }
    } else {
        None
    };

    (
        Status::Ok,
        Json(JudgeVoteProjectInfo::new(
            curr_proj.clone().map(|p| p.name),
            curr_proj.map(|p| p.location),
            prev_proj.clone().map(|p| p.name),
            prev_proj.map(|p| p.location),
        )),
    )
}

#[rocket::get("/judge/projects")]
pub async fn get_judge_seen_projects(
    db: &State<Arc<Database>>,
    token: Token,
) -> (Status, Json<Vec<Project>>) {
    // Get judge by token
    let judge = match find_judge_by_token(db, &token.0).await {
        Ok(j) => j,
        Err(e) => return (e, Json(Vec::<Project>::new())),
    };
    
    match find_projects_by_id(db, judge.seen_projects).await {
        Ok(p) => (Status::Ok, Json(p)),
        Err(e) => {
            eprintln!("{}", e);
            (Status::InternalServerError, Json(Vec::<Project>::new()))
        }

    }
}