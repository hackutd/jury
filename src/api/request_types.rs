use serde::Deserialize;

use crate::db::models::{Project, Judge};

#[derive(Deserialize)]
#[serde()]
pub struct Login<'r> {
    pub code: &'r str,
}

#[derive(Deserialize)]
#[serde()]
pub struct NewJudge {
    pub name: String,
    pub email: String,
    pub notes: Option<String>,
}

impl Into<Judge> for NewJudge {
    fn into(self) -> Judge {
        Judge::new(self.name, self.email, self.notes.unwrap_or("".to_string()))
    }
}

#[derive(Deserialize)]
#[serde()]
pub struct NewProject {
    pub name: String,
    pub description: String,
    pub try_link: Option<String>,
    pub video_link: Option<String>,
    pub challenge_list: Vec<String>,
}

impl Into<Project> for NewProject {
    fn into(self) -> Project {
        Project::new(
            self.name,
            self.description,
            self.try_link,
            self.video_link,
            self.challenge_list,
        )
    }
}