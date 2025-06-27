use std::net::SocketAddr;
use std::time::Duration;

use axum::Router;
use dotenvy::dotenv;
use hyper::server::conn::http1;
use hyper::service::service_fn;
use hyper_util::rt::tokio::TokioIo;
use reqwest::Client;
use serde_json::json;
use tokio::net::TcpListener;
use tower_service::Service;
use tracing::{error, info};
use tracing_subscriber::FmtSubscriber;
use tower_http::cors::{CorsLayer, Any};

mod api;
mod ollama;

#[tokio::main]
async fn main() {
    dotenv().ok();
    FmtSubscriber::builder().with_env_filter("info").init();

    let cors = CorsLayer::new()
        .allow_origin(Any)
        .allow_methods(Any)
        .allow_headers(Any);

    let app = api::routes().layer(cors);

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
