use mongodb::{bson::doc, options::ClientOptions, Client, Database};
use std::env;

pub async fn init_db() -> mongodb::error::Result<Database> {
    // Load MongoDB connection string from .env
    let client_uri = env::var("MONGODB_URI").expect("MONGODB_URI not defined");

    let options = ClientOptions::parse(&client_uri).await?;
    let client = Client::with_options(options)?;
    let db = client.database("jury");
    println!("Connecting to database...");

    // Ping the server to see if you can connect to the cluster
    client
        .database("admin")
        .run_command(doc! {"ping": 1}, None)
        .await?;
    println!("Connected successfully.");

    Ok(db)
}

// Initialize admin collection in MongoDB with one document lol
// This will represent all the admin settings
pub async fn admin_settings() {}
