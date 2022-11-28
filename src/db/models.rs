use mongodb::bson::oid::ObjectId;
use serde::{Deserialize, Serialize};
use chrono::{Utc, DateTime};

#[derive(Deserialize, Serialize, Debug)]
pub struct Project {
    #[serde(rename = "_id", skip_serializing_if = "Option::is_none")]
    id: Option<ObjectId>,
    name: String,
    location: String,
    description: String,
    mu: f64,
    sigma_sq: f64,
    active: bool,
    prioritized: bool,
    #[serde(with = "bson::serde_helpers::chrono_datetime_as_bson_datetime")]
    last_activity: DateTime<Utc>,
}

#[derive(Deserialize, Serialize, Debug)]
pub struct Judge {
    #[serde(rename = "_id", skip_serializing_if = "Option::is_none")]
    id: Option<ObjectId>,
    name: String,
    email: String,
    active: bool,
    #[serde(with = "bson::serde_helpers::chrono_datetime_as_bson_datetime")]
    last_activity: DateTime<Utc>,
    read_welcome: bool,
    notes: String,
    secret: String,
    next: ObjectId,
    prev: ObjectId,
    alpha: f64,
    beta: f64,
}

#[derive(Deserialize, Serialize, Debug)]
pub struct View {
    #[serde(rename = "_id", skip_serializing_if = "Option::is_none")]
    id: Option<ObjectId>,
    project: ObjectId,
    judge: ObjectId,
}

#[derive(Deserialize, Serialize, Debug)]
pub struct Decision {
    #[serde(rename = "_id", skip_serializing_if = "Option::is_none")]
    id: Option<ObjectId>,
    judge: ObjectId,
    winner: ObjectId,
    loser: ObjectId,
    #[serde(with = "bson::serde_helpers::chrono_datetime_as_bson_datetime")]
    time: DateTime<Utc>,
}

#[derive(Deserialize, Serialize, Debug)]
pub struct Flag {
    #[serde(rename = "_id", skip_serializing_if = "Option::is_none")]
    id: Option<ObjectId>,
    judge: ObjectId,
    project: ObjectId,
    reason: String,
    resolved: bool,
    #[serde(with = "bson::serde_helpers::chrono_datetime_as_bson_datetime")]
    time: DateTime<Utc>,
}
