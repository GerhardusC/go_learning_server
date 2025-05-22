use crate::ARGS;
use redis::Commands;
use chrono::Local;
use std::time::Instant;
use redis;
use rusqlite::{params, Connection};
use color_eyre::Result;

enum SensorValue {
    NUM(f32),
    STR(String)
}

pub struct MosquittoMessage {
    timestamp: u64,
    topic: String,
    value: SensorValue,
}

impl MosquittoMessage {
    fn new(timestamp: u64, topic: String, value: String) -> MosquittoMessage {
        let parsed_val = if let Ok(float) = value.parse::<f32>() {
            SensorValue::NUM(float)
        } else {
            SensorValue::STR(value)
        };
        MosquittoMessage {
            timestamp,
            topic,
            value: parsed_val,
        }
    }

    fn add_to_db(&self) -> Result<()> {
        let connection = Connection::open(&ARGS.db_path)?;
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

fn save_value_to_redis (topic: String, msg: String) -> Result<()> {
    let performance_now = Instant::now();
    let redis_client = redis::Client::open("redis://127.0.0.1/")?;
    let mut connection = redis_client.get_connection()?;
    let res = connection.set::<&str, &str, String>(&topic, &msg);
    match res {
        Ok(_) => {
            println!(
                "\x1b[1;30;104m{}:\x1b[0m \x1b[1m{}\x1b[0m added in \x1b[1;30;104mREDIS\x1b[0m at {}",
                topic,
                msg,
                Local::now().timestamp()
            );
            println!(
                "\x1b[1;30;47mSave to \x1b[0m\x1b[1;30;106mREDIS\x1b[0m duration: \x1b[1;96m{}\x1b[0m us",
                &performance_now.elapsed().as_micros()
            )
        },
        Err(err) => {
            println!("\x1b[1;30;41mError:\x1b[0m Something went wrong while adding data to redis: {}", err)
        },
    }
    Ok(())
}

pub fn save_value_to_redis_bg (topic: String, msg: String) {
    tokio::spawn(async move {
        match save_value_to_redis(topic, msg) {
            Ok(_) => (),
            Err(err) => {
                    println!("\x1b[1;30;41mError:\x1b[0m Something went wrong while connecting to redis: {}", err)
            },
        }
    });
}

pub fn save_value_to_db_bg (topic: String, msg: String) {
    tokio::spawn(async move {
        let performance_now = Instant::now();
        let now = Local::now().timestamp();
        let reading = MosquittoMessage::new(
            now.try_into().unwrap_or(now as u64),
            topic,
            msg,
        );
        match reading.add_to_db() {
            Ok(_) => {
                println!(
                    "\x1b[1;30;47mSave to \x1b[0m\x1b[1;30;45mDB\x1b[0m duration: \x1b[1;35m{}\x1b[0m us",
                    &performance_now.elapsed().as_micros()
                )
            },
            Err(err) => {
                println!("Something went wrong while adding data to the db: {}", err)
            },
        }
    });
}

