// src/state.rs
use std::sync::Arc;
use sqlx::PgPool;
use redis::Client;

#[derive(Clone)]
pub struct AppState {
    pub db: PgPool,
    pub auth_config: Arc<crate::auth::middleware::AuthConfig>,
    pub redis: Arc<Client>,
}
