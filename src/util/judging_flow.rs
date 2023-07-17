use std::error::Error;

use bson::oid::ObjectId;
use mongodb::Database;
use rand::{seq::SliceRandom, thread_rng, Rng};

use crate::db::{
    models::{Judge, Project},
    project::{find_all_active_projects, find_all_busy_projects, find_project_by_id},
};

use super::{
    constants::MIN_VIEWS,
    crowd_bt::{expected_information_gain, EPSILON},
};

pub async fn pick_next_project(db: &Database, judge: &Judge) -> Result<Option<Project>, Box<dyn Error>> {
    // Get items and shuffle them
    let mut items = find_preferred_items(db).await?;
    items.shuffle(&mut thread_rng());

    // If there are no items, return an error
    if items.is_empty() {
        return Ok(None)
    }

    // Randomly pick the first element with an EPSILON probability
    // Otherwise return the item with the highest expected information gain
    if thread_rng().gen::<f64>() < EPSILON {
        return Ok(Some(items[0].clone()));
    }

    // Get the judge's previous project
    let prev = if judge.prev.is_none() {
        Project::default()
    } else {
        let prev_id = judge
            .prev
            .ok_or::<Box<dyn Error>>("Invalid previous id".into())?;

        find_project_by_id(db, prev_id).await?
    };

    // Find max expected info gain
    Ok(Some(max_info_gain(items, &prev, judge)))
}

/// Find all projects that are higher priority with the following heuristic:
///  1. Ignore all projects that are inactive
///  2. If there are prioritized projects, pick from that list
///  3. If there are projects not currently being judged, pick from that list
///  4. If there are projects that have less than MIN_VIEWS, pick from that list
///  5. Return the remaining list, filtering at steps 2-4 if applicable
pub async fn find_preferred_items(db: &Database) -> Result<Vec<Project>, Box<dyn Error>> {
    // Get the list of all active projects
    let mut projects = find_all_active_projects(db).await?;

    // If there are no projects, return an empty list
    if projects.is_empty() {
        return Ok(projects);
    }

    // If there are prioritized projects, filter out the non-prioritized ones
    if projects.iter().any(|p| p.prioritized) {
        projects = projects.into_iter().filter(|p| p.prioritized).collect();
    }

    // Get all projects that are currently being judged
    // And remove busy projects from the list of projects
    let busy_project_ids = find_all_busy_projects(db).await?;
    projects = projects
        .into_iter()
        .filter(|p| !busy_project_ids.contains(&p.id.unwrap_or(ObjectId::new())))
        .collect();

    // If any projects have less than MIN_VIEWS, filter out all projects that have more than MIN_VIEWS
    if projects.iter().any(|p| p.seen < MIN_VIEWS) {
        projects = projects
            .into_iter()
            .filter(|p| p.seen < MIN_VIEWS)
            .collect();
    }

    Ok(projects)
}

pub fn max_info_gain(list: Vec<Project>, prev: &Project, judge: &Judge) -> Project {
    // Iterate through the list of projects and calculate the information gain for each
    let info_gain = list
        .iter()
        .map(|p| {
            expected_information_gain(
                judge.alpha,
                judge.beta,
                prev.mu,
                prev.sigma_sq,
                p.mu,
                p.sigma_sq,
            )
        })
        .collect::<Vec<f64>>();

    let max_index = info_gain
        .into_iter()
        .enumerate()
        .max_by(|(_, a), (_, b)| a.total_cmp(b))
        .map(|(index, _)| index)
        .unwrap_or(0);

    list[max_index].clone()
}
