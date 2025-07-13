use std::net::SocketAddr;
use std::sync::Arc;

use dotenvy::dotenv;
use reqwest::Client;
use tokio::net::TcpListener;
use tower_http::cors::{Any, CorsLayer};
use tracing::{info, error};
use tracing_subscriber::FmtSubscriber;

mod api;
mod ollama;
mod auth;
mod consts;
mod db;
mod state;
mod tables;

use crate::auth::middleware::{user_auth_middleware, AuthConfig};
use crate::consts::user_authentication_endpoint;
use crate::state::AppState;

#[tokio::main]
async fn main() {
    dotenv().ok();

    FmtSubscriber::builder().with_env_filter("info").init();

    let pool = db::connect().await;

    let cors = CorsLayer::new()
        .allow_origin(Any)
        .allow_methods(Any)
        .allow_headers(Any);

    let auth_config = Arc::new(AuthConfig {
        auth_service_url: user_authentication_endpoint(),
        http_client: Client::new(),
    });

    let shared_state = AppState {
        db: pool,
        auth_config,
    };

    let app = api::routes()
        .layer(cors)
        .layer(axum::middleware::from_fn_with_state(
            shared_state.clone(),
            user_auth_middleware,
        ))
        .with_state(shared_state); // ✅ this is the key

    let addr = SocketAddr::from(([0, 0, 0, 0], 3000));
    info!("🚀 Server running at http://{addr}");

    let listener = TcpListener::bind(addr)
        .await
        .expect("Failed to bind TCP listener");

    // Now this works:
    if let Err(e) = axum::serve(listener, app.into_make_service()).await {
        error!("❌ Server error: {e}");
    }
}
