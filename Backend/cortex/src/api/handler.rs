use axum::{
    extract::{Query, State},
    response::Response,
    routing::get,
    Router,
    body::Body,
    Extension,
};
use http_body_util::StreamBody;
use std::{sync::{Arc, Mutex}, convert::Infallible};
use uuid::Uuid;
use sqlx::PgPool;
use futures_util::TryStreamExt;
use crate::ollama::client::query_ollama_stream;
use crate::auth::types::UserObject;
use crate::tables::{Session, Message};

#[derive(Debug, serde::Deserialize)]
pub struct AskParams {
    pub q: String,
    pub session_id: Option<Uuid>,
    pub algorithm_id: Option<Uuid>,
}

pub fn routes() -> Router<Arc<PgPool>> {
    Router::new().route("/v1/ask", get(handle_ask))
}

async fn handle_ask(
    Query(params): Query<AskParams>,
    Extension(user): Extension<UserObject>,
    State(pool): State<Arc<PgPool>>,
) -> Result<Response, Infallible> {
    let prompt = params.q.clone();
    let algorithm_id = params.algorithm_id;
    let user_id = user.id;

    // Create a new session if session_id not provided
    let session_id = if let Some(sid) = params.session_id {
        sid
    } else {
        match sqlx::query_as::<_, Session>(
            r#"
            INSERT INTO sessions (user_id, algorithm_id, title)
            VALUES ($1, $2, $3)
            RETURNING id, created_at, user_id, algorithm_id, title
            "#
        )
        .bind(user_id)
        .bind(algorithm_id)
        .bind(&prompt)
        .fetch_one(&*pool)
        .await
        {
            Ok(session) => session.id,
            Err(e) => {
                let err_msg = format!("❌ Failed to create session: {e}");
                return Ok(Response::builder()
                    .status(500)
                    .header("content-type", "text/plain")
                    .body(err_msg.into())
                    .unwrap());
            }
        }
    };

    // Set up thread-safe system message collector
    let system_msg = Arc::new(Mutex::new(String::new()));
    let system_msg_clone = Arc::clone(&system_msg);

    // Stream the LLM response and collect it in background
    let stream = match query_ollama_stream(prompt.clone()).await {
        Ok(stream) => stream.inspect_ok(move |chunk| {
            if let Ok(text) = std::str::from_utf8(chunk) {
                if let Ok(mut lock) = system_msg_clone.lock() {
                    lock.push_str(text);
                }
            }
        }),
        Err(e) => {
            let msg = format!("❌ LLM Error: {}\n", e);
            return Ok(Response::builder()
                .status(500)
                .header("content-type", "text/plain")
                .body(msg.into())
                .unwrap());
        }
    };

    // Return streaming response to client
    let body = Body::from_stream(StreamBody::new(stream));
    let response = Response::builder()
        .status(200)
        .header("content-type", "text/event-stream")
        .body(body)
        .unwrap();

    // Clone things needed in background
    let pool = Arc::clone(&pool);
    let prompt_clone = prompt.clone();

    tokio::spawn(async move {
        // Extract collected message from Arc<Mutex<_>>
        let collected = Arc::try_unwrap(system_msg)
            .ok()
            .and_then(|m| m.into_inner().ok())
            .unwrap_or_default();

        if let Err(e) = sqlx::query_as::<_, Message>(r#"
            INSERT INTO messages (session_id, user_message, system_message)
            VALUES ($1, $2, $3)
            RETURNING id, created_at, session_id, user_message, system_message
        "#)
        .bind(session_id)
        .bind(prompt_clone)
        .bind(collected)
        .fetch_one(&*pool)
        .await
        {
            tracing::error!("❌ Failed to insert message: {}", e);
        }
    });

    Ok(response)
}
