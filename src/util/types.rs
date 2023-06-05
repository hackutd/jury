use std::time::SystemTime;

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
    pub avg_votes: f64,
    pub num_active: u64,
}

impl Default for JudgeStats {
    fn default() -> Self {
        Self {
            num: 0,
            avg_votes: 0.0,
            num_active: 0,
        }
    }
}

#[derive(Serialize, Deserialize)]
#[serde()]
pub struct ProjectStats {
    pub num: u64,
    pub avg_votes: f64,
    pub avg_seen: f64,
}

impl Default for ProjectStats {
    fn default() -> Self {
        Self {
            num: 0,
            avg_votes: 0.0,
            avg_seen: 0.0,
        }
    }
}

#[derive(FromForm)]
pub struct CsvUpload {
    pub csv: String,

    #[field(name = uncased("hasHeader"), default = false)]
    pub has_header: bool,
}

#[derive(Serialize)]
pub struct ClockState {
    pub start: u128,
    pub prev: u128,
    pub paused: bool,
}

impl ClockState {
    pub fn new() -> Self {
        Self {
            start: 0,
            prev: 0,
            paused: true,
        }
    }

    pub fn pause(&mut self) {
        if self.paused {
            return;
        }
        self.paused = true;
        self.prev += Self::get_curr_time() - self.start;
    }

    pub fn resume(&mut self) {
        if !self.paused {
            return;
        }
        self.paused = false;
        self.start = Self::get_curr_time();
    }

    pub fn reset(&mut self) {
        self.start = 0;
        self.prev = 0;
        self.paused = false;
    }

    pub fn get_duration(&self) -> u128 {
        if self.paused {
            self.prev
        } else {
            self.prev + (Self::get_curr_time() - self.start)
        }
    }

    pub fn get_curr_time() -> u128 {
        SystemTime::now()
            .duration_since(SystemTime::UNIX_EPOCH)
            .expect("Somehow the current time is before the UNIX epoch :thinking emoji:")
            .as_millis()
    }

    pub fn get_copy(&self) -> Self {
        Self {
            start: self.start,
            prev: self.prev,
            paused: self.paused,
        }
    }
}
