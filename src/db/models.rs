use bson::doc;
use chrono::{DateTime, TimeZone, Utc};
use mongodb::{bson::oid::ObjectId, Collection, Database};
use rand::Rng;
use serde::{Deserialize, Serialize};

use crate::util::crowd_bt;

#[derive(Deserialize, Serialize, Debug)]
pub struct Options {
    #[serde(rename = "_id", skip_serializing_if = "Option::is_none")]
    pub id: Option<ObjectId>,
    pub next_table_num: u32,
}

impl Options {
    pub fn new() -> Self {
        Self {
            id: None,
            next_table_num: 1,
        }
    }

    /// Increment next table num and return the value
    pub fn get_next_table_num(&mut self) -> u32 {
        self.next_table_num += 1;
        self.next_table_num - 1
    }

    /// Saves the options to the database
    pub async fn save(&self, db: &Database) -> Result<(), Box<dyn std::error::Error>> {
        let collection: Collection<Options> = db.collection("options");
        collection
            .update_one(
                doc! {},
                doc! {
                    "$set": {
                        "next_table_num": self.next_table_num,
                    }
                },
                None,
            )
            .await?;
        Ok(())
    }
}

impl Clone for Options {
    fn clone(&self) -> Self {
        Options {
            id: self.id.clone(),
            next_table_num: self.next_table_num,
        }
    }
}

#[derive(Deserialize, Serialize, Debug)]
pub struct Project {
    #[serde(rename = "_id", skip_serializing_if = "Option::is_none")]
    pub id: Option<ObjectId>,
    pub name: String,
    pub location: u64,
    pub description: String,
    pub try_link: Option<String>,
    pub video_link: Option<String>,
    pub challenge_list: Vec<String>,
    pub seen: u64,
    pub votes: u64,
    pub mu: f64,
    pub sigma_sq: f64,
    pub active: bool,
    pub prioritized: bool,
    #[serde(with = "bson::serde_helpers::chrono_datetime_as_bson_datetime")]
    pub last_activity: DateTime<Utc>,
}

impl Project {
    pub fn new(
        name: String,
        description: String,
        try_link: Option<String>,
        video_link: Option<String>,
        challenge_list: Vec<String>,
    ) -> Self {
        Self {
            id: None,
            name,
            location: 0,
            description,
            try_link,
            video_link,
            challenge_list,
            seen: 0,
            votes: 0,
            mu: 0f64,
            sigma_sq: 1f64,
            active: true,
            prioritized: false,
            last_activity: Utc.timestamp_opt(0, 0).unwrap(),
        }
    }
}

#[derive(Deserialize, Serialize, Debug)]
pub struct Judge {
    #[serde(rename = "_id", skip_serializing_if = "Option::is_none")]
    pub id: Option<ObjectId>,
    pub token: String,
    pub code: String,
    pub name: String,
    pub email: String,
    pub active: bool,
    #[serde(with = "bson::serde_helpers::chrono_datetime_as_bson_datetime")]
    pub last_activity: DateTime<Utc>,
    pub read_welcome: bool,
    pub notes: String,
    pub votes: u64,
    pub next: Option<ObjectId>,
    pub prev: Option<ObjectId>,
    pub alpha: f64,
    pub beta: f64,
}

impl Judge {
    pub fn new(name: String, email: String, notes: String) -> Self {
        Self {
            id: None,
            token: "".to_string(),
            code: rand::thread_rng().gen_range(100000..999999).to_string(),
            name,
            email,
            active: true,
            read_welcome: false,
            notes,
            votes: 0,
            next: None,
            prev: None,
            alpha: crowd_bt::ALPHA_PRIOR,
            beta: crowd_bt::BETA_PRIOR,
            last_activity: Utc.timestamp_opt(0, 0).unwrap(),
        }
    }

    pub fn default() -> Self {
        Self::new("".to_string(), "".to_string(), "".to_string())
    }
}

#[derive(Deserialize, Serialize, Debug)]
pub struct View {
    #[serde(rename = "_id", skip_serializing_if = "Option::is_none")]
    pub id: Option<ObjectId>,
    pub project: ObjectId,
    pub judge: ObjectId,
}

#[derive(Deserialize, Serialize, Debug)]
pub struct Decision {
    #[serde(rename = "_id", skip_serializing_if = "Option::is_none")]
    pub id: Option<ObjectId>,
    pub judge: ObjectId,
    pub winner: ObjectId,
    pub loser: ObjectId,
    #[serde(with = "bson::serde_helpers::chrono_datetime_as_bson_datetime")]
    pub time: DateTime<Utc>,
}

#[derive(Deserialize, Serialize, Debug)]
pub struct Flag {
    #[serde(rename = "_id", skip_serializing_if = "Option::is_none")]
    pub id: Option<ObjectId>,
    pub judge: ObjectId,
    pub project: ObjectId,
    pub reason: String,
    pub resolved: bool,
    #[serde(with = "bson::serde_helpers::chrono_datetime_as_bson_datetime")]
    pub time: DateTime<Utc>,
}
