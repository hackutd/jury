use rocket::FromForm;
use serde::{Deserialize, Serialize};

#[derive(Deserialize)]
#[serde()]
pub struct AdminLogin<'r> {
    pub password: &'r str,
}

#[derive(Serialize, Deserialize)]
#[serde()]
pub struct Stats {
    pub projects: u64,
    pub seen: u64,
    pub votes: u64,
    pub time: u64,
    pub avg_mu: f64,
    pub avg_sigma: f64,
    pub judges: u64,
}

impl Default for Stats {
    fn default() -> Self {
        Self {
            projects: 0,
            seen: 0,
            votes: 0,
            time: 0,
            avg_mu: 0.0,
            avg_sigma: 0.0,
            judges: 0,
        }
    }
}

#[derive(Serialize, Deserialize)]
#[serde()]
pub struct JudgeStats {
    pub num: u64,
    pub alpha: f64,
    pub beta: f64,
}

impl Default for JudgeStats {
    fn default() -> Self {
        Self {
            num: 0,
            alpha: 0.0,
            beta: 0.0,
        }
    }
}

#[derive(FromForm)]
pub struct CsvUpload {
    pub csv: String,
    pub has_header: bool,
}
