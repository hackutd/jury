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
