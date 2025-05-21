use crate::ARGS;
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
    pub fn new(timestamp: u64, topic: String, value: String) -> MosquittoMessage {
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

    pub fn add_to_db(&self) -> Result<()> {
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

