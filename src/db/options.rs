use std::error::Error;

use bson::doc;
use mongodb::{Collection, Database};

use crate::db::models::Options;

pub async fn get_options(db: &Database) -> Result<Options, Box<dyn Error>> {
    let collection: Collection<Options> = db.collection("options");
    let content = collection.find_one(None, None).await?;
    match content {
        Some(x) => Ok(x),
        None => {
            let options = Options::new();
            collection.insert_one(options.clone(), None).await?;
            Ok(options)
        }
    }
}

pub async fn get_next_location(db: &Database) -> Result<u64, mongodb::error::Error> {
    let collection = db.collection::<Options>("options");
    let res = collection.find_one(None, None).await?;
    let options = match res {
        Some(o) => o,
        None => {
            let new_options = Options::new();
            collection.insert_one(new_options.clone(), None).await?;
            new_options
        }
    };
    Ok(options.next_table_num.into())
}

pub async fn increment_location(
    db: &Database,
    amt: Option<u32>,
) -> Result<(), mongodb::error::Error> {
    let collection = db.collection::<Options>("options");
    let res = collection.find_one(None, None).await?;
    let options = match res {
        Some(o) => o,
        None => {
            let new_options = Options::new();
            collection.insert_one(new_options.clone(), None).await?;
            new_options
        }
    };

    let amt = amt.unwrap_or(1);
    collection
        .update_one(
            doc! { "_id": options.id },
            doc! { "$set": { "next_table_num": options.next_table_num + amt } },
            None,
        )
        .await?;
    Ok(())
}
