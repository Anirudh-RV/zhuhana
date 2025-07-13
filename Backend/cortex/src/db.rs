use sqlx::{migrate::Migrator, PgPool, postgres::PgPoolOptions};
use std::{env};

static MIGRATOR: Migrator = sqlx::migrate!("./migrations");

pub async fn connect() -> PgPool {
    let host = env::var("DB_HOST").expect("DB_HOST not set");
    let port = env::var("DB_PORT").expect("DB_PORT not set");
    let user = env::var("DB_USER").expect("DB_USER not set");
    let pass = env::var("DB_PASSWORD").expect("DB_PASSWORD not set");
    let db = env::var("DB_NAME").expect("DB_NAME not set");

    let url = format!("postgres://{user}:{pass}@{host}:{port}/{db}");

    let pool = PgPoolOptions::new()
        .max_connections(5)
        .connect(&url)
        .await
        .expect("❌ Failed to connect to DB");

    // Run pending migrations automatically
    MIGRATOR
        .run(&pool)
        .await
        .expect("❌ Failed to run DB migrations");

    pool
}
