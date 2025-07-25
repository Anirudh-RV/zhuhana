use redis::Client;
use std::env;
use std::error::Error;

/// Initializes and returns a Redis client using env variables.
pub fn init_redis_client() -> Result<Client, Box<dyn Error>> {
    let redis_url = build_redis_url()?;
    let client = Client::open(redis_url)?;
    Ok(client)
}

fn build_redis_url() -> Result<String, env::VarError> {
    let host = env::var("REDIS_HOST")?;
    let port = env::var("REDIS_PORT")?;
    let password = env::var("REDIS_PASSWORD")?;

    Ok(format!("redis://:{}@{}:{}/", password, host, port))
}
