use crate::utils::validate_topic;
use crate::ARGS;
use crate::data_collection::{save_value_to_db_bg, save_value_to_redis_bg};
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
    let mosquitto_client = Client::with_auto_id()?;
    mosquitto_client
        .connect(&ARGS.broker_ip, 1883, Duration::from_secs(5), None)
        .await?;

    let subscriptions = mosquitto_client.subscriber();
    let topic = validate_topic(&ARGS.base_topic);
    println!("\x1b[1;30;43mINFO:\x1b[47m Base topic:\x1b[0m {}", topic);

    mosquitto_client.subscribe(&topic, QoS::AtMostOnce).await?;

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
            let msg_str = String::from_utf8(message.payload)?;

            save_value_to_db_bg(message.topic.clone(), msg_str.clone());
            save_value_to_redis_bg(message.topic, msg_str);
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

