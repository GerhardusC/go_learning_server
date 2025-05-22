pub fn validate_topic(topic: &str) -> String {
    if topic == "" {
        return "/#".to_owned()
    }

    if topic == "/#" {
        return "/#".to_owned()
    }

    let topic_as_bytes = topic.as_bytes();
    for (i, char) in topic_as_bytes.iter().enumerate() {
        if let Some(next_char) = topic_as_bytes.get(i + 1) {
            if *char == *next_char && *char == 47 {
                println!("\x1b[1;30;41mError:\x1b[0m cannot have two consecutive slashes in topic. falling back to /#");
                return "/#".to_owned()
            }
        }
    }

    if topic.starts_with("/") && topic.ends_with("/") {
        return format!("{}#", topic)
    }

    if !topic.starts_with("/") && topic.ends_with("/") {
        return format!("/{}#", topic)
    }

    if topic.starts_with("/") && !topic.ends_with("/") {
        if topic.ends_with("/#") {
            return format!("{}", topic)
        } else {
            return format!("{}/#", topic)
        }
    }

    return format!("/{}/#", topic)
}
