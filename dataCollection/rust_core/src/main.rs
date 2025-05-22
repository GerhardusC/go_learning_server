mod args;
mod data_collection;
mod db_setup;
mod event_handling;
mod utils;

use args::ARGS;
use event_handling::start_subscription_loop;
use db_setup::setup_db;

use color_eyre::Result;


#[tokio::main]
async fn main() -> Result<()> {
    color_eyre::install()?;
    setup_db()?;
    start_subscription_loop().await;
    Ok(())
}

