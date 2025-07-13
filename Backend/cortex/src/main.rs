use std::net::SocketAddr;
use std::sync::Arc;

use axum::Extension;
use dotenvy::dotenv;
use hyper::server::conn::http1;
use hyper::service::service_fn;
use hyper_util::rt::tokio::TokioIo;
use reqwest::Client;
use tokio::net::TcpListener;
use tower_service::Service;
use tower::ServiceBuilder;
use tracing::{error, info};
use tracing_subscriber::FmtSubscriber;
use tower_http::cors::{CorsLayer, Any};


mod api;
mod ollama;
mod auth;
mod consts;


use crate::auth::middleware::{user_auth_middleware, AuthConfig};
use crate::consts::user_authentication_endpoint;


#[tokio::main]
async fn main() {
    dotenv().ok();
    FmtSubscriber::builder().with_env_filter("info").init();

    let cors = CorsLayer::new()
        .allow_origin(Any)
        .allow_methods(Any)
        .allow_headers(Any);

    let auth_config = Arc::new(AuthConfig {
        auth_service_url: user_authentication_endpoint(),
        http_client: Client::new(),
    });

    let app = api::routes().layer(
        ServiceBuilder::new()
            .layer(Extension(auth_config))
            .layer(axum::middleware::from_fn(user_auth_middleware))
            .layer(cors),
    );

    let addr = SocketAddr::from(([0, 0, 0, 0], 3000));
    info!("🚀 Running on http://{addr}");

    let listener = TcpListener::bind(addr)
        .await
        .expect("Failed to bind TCP listener");

    loop {
        let (stream, _) = listener.accept().await.expect("Failed to accept connection");
        let app = app.clone();

        tokio::spawn(async move {
            let io = TokioIo::new(stream);
            if let Err(err) = http1::Builder::new()
                .serve_connection(io, service_fn(move |req| app.clone().call(req)))
                .await
            {
                error!("Server error: {err}");
            }
        });
    }
}
