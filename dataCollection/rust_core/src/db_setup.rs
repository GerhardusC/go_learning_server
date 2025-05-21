use rusqlite::Connection;
use color_eyre::Result;
use crate::ARGS;

pub fn setup_db () -> Result<()> {
    let connection = Connection::open(&ARGS.db_path)?;
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
                value varchar(255)
        )
        ",
        (),
    )?;
    Ok(())
}

