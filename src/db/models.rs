use mongodb::bson::oid::ObjectId;
use serde::{Deserialize, Serialize};
use chrono::{Utc, DateTime};

#[derive(Deserialize, Serialize, Debug)]
pub struct Project {
    #[serde(rename = "_id", skip_serializing_if = "Option::is_none")]
    pub id: Option<ObjectId>,
    pub name: String,
    pub location: String,
    pub description: String,
    pub seen: u64,
    pub votes: u64,
    pub mu: f64,
    pub sigma_sq: f64,
    pub active: bool,
    pub prioritized: bool,
    #[serde(with = "bson::serde_helpers::chrono_datetime_as_bson_datetime")]
    pub last_activity: DateTime<Utc>,
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
    pub next: Option<ObjectId>,
    pub prev: Option<ObjectId>,
    pub alpha: f64,
    pub beta: f64,
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
