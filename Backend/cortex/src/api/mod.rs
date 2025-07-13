use axum::{
    extract::Query,
    response::Response,
    routing::get,
    Router,
};
use axum::Extension;
use axum::body::Body;
use http_body_util::StreamBody;
use std::{collections::HashMap, convert::Infallible};

use crate::ollama::client::query_ollama_stream;
use crate::auth::types::UserObject;

pub fn routes() -> Router {
    Router::new().route("/v1/ask", get(handle_ask))
}

async fn handle_ask(
    Query(params): Query<HashMap<String, String>>,
    Extension(user): Extension<UserObject>,
) -> Result<Response, Infallible> {
    let prompt = match params.get("q") {
        Some(p) => p.clone(),
        None => {
            let msg = "Missing `q` parameter\n";
            return Ok(Response::builder()
                .status(400)
                .header("content-type", "text/plain")
                .body(msg.into())
                .unwrap());
        }
    };

    match query_ollama_stream(prompt).await {
        Ok(stream) => {
            let body = Body::from_stream(StreamBody::new(stream));
            let response = Response::builder()
                .status(200)
                .header("content-type", "text/event-stream")
                .body(body)
                .unwrap();

            Ok(response)
        }
        Err(e) => {
            let msg = format!("Error: {}\n", e);
            Ok(Response::builder()
                .status(500)
                .header("content-type", "text/plain")
                .body(msg.into())
                .unwrap())
        }
    }
}
