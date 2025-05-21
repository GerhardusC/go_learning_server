use clap::Parser;
use std::sync::LazyLock;

/// Lightweight program that performs data collection via MQTT and saves the data to a SQLITE database.
#[derive(Parser, Debug, Clone)]
#[command(version, about, long_about = None)]
pub struct Cli {
    /// Path to database
    #[arg(short, long, default_value_t = String::from("./dev.db"))]
    pub db_path: String,
    /// Base topic to subscribe to, if omitted, you will be subscribed to /#
    #[arg(short('t'), long, default_value_t = String::from(""))]
    pub base_topic: String,
    /// IP or domain name of mqtt broker
    #[arg(short, long, default_value_t = String::from("localhost"))]
    pub broker_ip: String,
}

pub static ARGS: LazyLock<Cli> = LazyLock::new(|| Cli::parse());
