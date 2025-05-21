use chrono::Local;
use clap::Parser;
use rusqlite::{params, Connection};
use std::{io::{Error, ErrorKind}, time::Duration};

use color_eyre::Result;
use mosquitto_rs::*;
use tokio::time::sleep;

/// Lightweight program that performs data collection via MQTT and saves the data to a SQLITE database.
#[derive(Parser, Debug, Clone)]
#[command(version, about, long_about = None)]
struct Cli {
    #[arg(short, long, default_value_t = String::from("./dev.db"))]
    db_path: String,

    #[arg(short('t'), long, default_value_t = String::from(""))]
    base_topic: String,

    #[arg(short, long, default_value_t = String::from("localhost"))]
    broker_ip: String,
}

enum SensorValue {
    NUM(f32),
    STR(String)
}

struct SensorReading {
    timestamp: u64,
    topic: String,
    value: SensorValue,
}

impl SensorReading {
    fn new(timestamp: u64, topic: String, value: String) -> SensorReading {
        let parsed_val = if let Ok(float) = value.parse::<f32>() {
            SensorValue::NUM(float)
        } else {
            SensorValue::STR(value)
        };
        SensorReading {
            timestamp,
            topic,
            value: parsed_val,
        }
    }

    fn add_to_db(&self, args: &Cli) -> Result<()> {
        let connection = Connection::open(&args.db_path)?;
        match &self.value {
            SensorValue::NUM(float_val) => {
                connection.execute(
                    "INSERT INTO MEASUREMENTS VALUES (?1, ?2, ?3);",
                    params![self.timestamp, self.topic, float_val],
                )?;
                println!(
                    "\x1b[1;30;104m{}:\x1b[0m\x1b[1m {}\x1b[0m saved to \x1b[1;30;46mMEASUREMENTS\x1b[0m at {}",
                    self.topic,
                    float_val,
                    self.timestamp
                )
            },
            SensorValue::STR(string_val) => {
                connection.execute(
                    "INSERT INTO LOGS VALUES (?1, ?2, ?3);",
                    params![self.timestamp, self.topic, string_val],
                )?;
                println!(
                    "\x1b[1;30;104m{}:\x1b[0m\x1b[1m {}\x1b[0m saved to \x1b[1;30;103mLOGS\x1b[0m at {}",
                    self.topic,
                    string_val,
                    self.timestamp
                )
            },
        }

        Ok(())
    }
}

#[tokio::main]
async fn main() -> Result<()> {
    color_eyre::install()?;
    let args = Cli::parse();

    let connection = Connection::open(&args.db_path)?;
    connection.execute("
        CREATE TABLE if not exists MEASUREMENTS (
                timestamp int,
                topic varchar(255),
                value float
        )
        ",
        (),
    )?;

    connection.execute("
        CREATE TABLE if not exists LOGS (
                timestamp int,
                topic varchar(255),
                value float
        )
        ",
        (),
    )?;

    loop {
        let success = subscribe(&args).await;
        match success {
            Ok(()) => (),
            Err(e) => {
                println!("Something went wrong, trying again soon.\n{}", e);
                sleep(Duration::from_secs(1)).await;
            }
        }
    }
}

async fn subscribe(args: &Cli) -> Result<()> {
    let client = Client::with_auto_id()?;
    client
        .connect(&args.broker_ip, 1883, Duration::from_secs(5), None)
        .await?;

    let subscriptions = client.subscriber();
    let topic = if args.base_topic == ""
        {"/#"} else
        {&format!("/{}/#", &args.base_topic)};

    client.subscribe(topic, QoS::AtMostOnce).await?;

    loop {
        if let Some(sub) = &subscriptions {
            match sub.recv().await {
                Ok(msg) => {
                    respond_to_event(msg, args)?;
                },
                Err(err) => {
                    println!("Error receiveing: {}", err)
                },
            }
        } else {
                println!("Error with event")
        }
    }
}

fn respond_to_event (msg: Event, args: &Cli) -> Result<()> {
    match msg {
        Event::Message(message) => {
            let now_timestamp_i64 = Local::now().timestamp();
            let now = if let Ok(val) = u64::try_from(now_timestamp_i64) {
                val
            } else {
                println!("Failed to convert u64 to i64");
                now_timestamp_i64 as u64
            };
            let reading = SensorReading::new(
                now,
                message.topic,
                String::from_utf8(message.payload)?,
            );

            reading.add_to_db(&args)?;
        },
        Event::Connected(connection_status) => {
            println!("\x1b[1;30;42mConnected:\x1b[0m {}", connection_status)
        },
        Event::Disconnected(reason_code) => {
            println!("\x1b[1;30;41mDisconnected:\x1b[0m {}", reason_code);
            return Err(
                Error::new(
                    ErrorKind::Other, "Something went wrong, disconnected."
                ).into()
            )
        },
    }
    Ok(())
}
