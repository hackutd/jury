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
            let options = Options::default();
            collection.insert_one(options.clone(), None).await?;
            Ok(options)
        }
    }
}

impl Options {
    pub fn default() -> Options {
        Options { next_table_num: 1u32 }
    }

    /// Increment next table num and return the value
    pub fn get_next_table_num(&mut self) -> u32 {
        self.next_table_num += 1;
        self.next_table_num - 1
    }

    /// Saves the options to the database
    pub async fn save(&self, db: &Database) -> Result<(), Box<dyn Error>> {
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
            next_table_num: self.next_table_num,
        }
    }
}
