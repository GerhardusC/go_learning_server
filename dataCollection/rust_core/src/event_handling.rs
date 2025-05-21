use crate::ARGS;
use crate::data_collection::MosquittoMessage;
use chrono::Local;
use std::{io::{Error, ErrorKind}, time::Duration};
use mosquitto_rs::*;
use tokio::time::sleep;
use color_eyre::Result;

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
    let client = Client::with_auto_id()?;
    client
        .connect(&ARGS.broker_ip, 1883, Duration::from_secs(5), None)
        .await?;

    let subscriptions = client.subscriber();
    let topic = if &ARGS.base_topic == ""
        {"/#"} else
        {&format!("/{}/#", &ARGS.base_topic)};

    client.subscribe(topic, QoS::AtMostOnce).await?;

    loop {
        if let Some(sub) = &subscriptions {
            match sub.recv().await {
                Ok(msg) => {
                    respond_to_event(msg)?;
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

fn respond_to_event (msg: Event) -> Result<()> {
    match msg {
        Event::Message(message) => {
            let now_timestamp_i64 = Local::now().timestamp();
            let now = if let Ok(val) = u64::try_from(now_timestamp_i64) {
                val
            } else {
                println!("Failed to convert u64 to i64");
                now_timestamp_i64 as u64
            };
            let reading = MosquittoMessage::new(
                now,
                message.topic,
                String::from_utf8(message.payload)?,
            );

            tokio::spawn(async move {
                match reading.add_to_db() {
                    Ok(_) => (),
                    Err(err) => {
                        println!("Something went wrong while adding data to the db: {}", err)
                    },
                }
            });
            ()
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
