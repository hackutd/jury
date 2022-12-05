use mongodb::{
    options::{ClientOptions, ResolverConfig},
    Client, Database,
};
use std::error::Error;
use std::env;

pub async fn init_db() -> Result<Database, Box<dyn Error>> {
    // Load MongoDB connection string from .env
    let client_uri = env::var("MONGODB_URI").expect("MONGODB_URI not defined");

    let options =
        ClientOptions::parse_with_resolver_config(&client_uri, ResolverConfig::cloudflare())
            .await?;
    let client = Client::with_options(options)?;
    let db = client.database("gavel");

    Ok(db)
}

// Initialize admin collection in MongoDB with one document lol
// This will represent all the admin settings
pub async fn admin_settings() {

}
