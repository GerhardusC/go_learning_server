use clap::Parser;
use std::sync::LazyLock;

/// Lightweight program that performs data collection via MQTT and saves the data to a SQLITE database.
#[derive(Parser, Debug, Clone)]
#[command(version, about, long_about = None)]
pub struct Cli {
    #[arg(short, long, default_value_t = String::from("./dev.db"))]
    pub db_path: String,

    #[arg(short('t'), long, default_value_t = String::from(""))]
    pub base_topic: String,

    #[arg(short, long, default_value_t = String::from("localhost"))]
    pub broker_ip: String,
}

pub static ARGS: LazyLock<Cli> = LazyLock::new(|| Cli::parse());
