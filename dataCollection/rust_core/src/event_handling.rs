use crate::ARGS;
use crate::data_collection::MosquittoMessage;
use redis::Commands;
use chrono::Local;
use std::{io::{Error, ErrorKind}, time::{Duration, Instant}};
use mosquitto_rs::*;
use tokio::time::sleep;
use color_eyre::Result;
use redis;

pub async fn start_subscription_loop () {
    loop {
        let success = subscribe_to_base_topic().await;
        match success {
            Ok(()) => (),
            Err(e) => {
                println!("Something went wrong, trying again soon.\n{}", e);
                sleep(Duration::from_secs(1)).await;
            }
        }
    }
}

async fn subscribe_to_base_topic() -> Result<()> {
    let redis_client = redis::Client::open("redis://127.0.0.1/")?;

    let mosquitto_client = Client::with_auto_id()?;
    mosquitto_client
        .connect(&ARGS.broker_ip, 1883, Duration::from_secs(5), None)
        .await?;

    let subscriptions = mosquitto_client.subscriber();
    let topic = if &ARGS.base_topic == ""
        {"/#"} else
        {&format!("/{}/#", &ARGS.base_topic)};

    mosquitto_client.subscribe(topic, QoS::AtMostOnce).await?;

    loop {
        if let Some(sub) = &subscriptions {
            match sub.recv().await {
                Ok(msg) => {
                    respond_to_event(msg, &redis_client)?;
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

fn respond_to_event (msg: Event, redis_client: &redis::Client) -> Result<()> {
    match msg {
        Event::Message(message) => {
            let msg_str = String::from_utf8(message.payload)?;

            let now_timestamp_i64 = Local::now().timestamp();
            let now = if let Ok(val) = u64::try_from(now_timestamp_i64) {
                val
            } else {
                println!("Failed to convert u64 to i64");
                now_timestamp_i64 as u64
            };
            let reading = MosquittoMessage::new(
                now,
                (&message.topic).to_owned(),
                (&msg_str).to_owned(),
            );

            let performance_now = Instant::now();
            tokio::spawn(async move {
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

            let performance_now = Instant::now();
            let connection = redis_client.get_connection();

            match connection {
                Ok(mut conn) => {
                    tokio::spawn(async move {
                        let res = conn.set::<&str, &str, String>(&message.topic, &msg_str);
                        match res {
                            Ok(_) => {
                                println!(
                                    "\x1b[1;30;104m{}:\x1b[0m \x1b[1m{}\x1b[0m added in \x1b[1;30;104mREDIS\x1b[0m at {}",
                                    message.topic,
                                    msg_str,
                                    now
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
                    });
                },
                Err(err) => {
                        println!("\x1b[1;30;41mError:\x1b[0m Something went wrong while connecting to redis: {}", err)
                },
            }

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
