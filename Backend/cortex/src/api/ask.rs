use axum::{
    extract::{Query, State, Extension},
    response::Response,
    body::Body,
};
use futures_util::{StreamExt};
use http_body_util::StreamBody;
use uuid::Uuid;
use std::{convert::Infallible};
use async_stream::stream;
use bytes::Bytes;
use tokio::sync::mpsc;

use crate::{
    auth::types::UserObject,
    ollama::client::query_ollama_stream,
    tables::{Message, Session},
    state::AppState,
};

#[derive(Debug, serde::Deserialize)]
pub struct AskParams {
    pub q: String,
    pub session_id: Option<Uuid>,
}

pub async fn handle_ask(
    Query(params): Query<AskParams>,
    Extension(user): Extension<UserObject>,
    State(state): State<AppState>,
) -> Result<Response, Infallible> {
    let pool = &state.db;
    let prompt = params.q.clone();

    let session_id = match params.session_id {
        Some(sid) => sid,
        None => {
            return Ok(Response::builder()
                .status(400)
                .header("content-type", "text/plain")
                .body("❌ Missing required `session_id` query parameter".into())
                .unwrap());
        }
    };

    let session_result = sqlx::query_as::<_, Session>(
        r#"
        SELECT id, created_at, user_id, algorithm_id, title
        FROM sessions
        WHERE id = $1 AND user_id = $2
        "#,
    )
    .bind(session_id)
    .bind(user.id)
    .fetch_optional(pool)
    .await;

    match session_result {
        Ok(Some(_)) => {
            // ✅ Session found and belongs to user — continue
        }
        Ok(None) => {
            return Ok(Response::builder()
                .status(403)
                .header("content-type", "text/plain")
                .body("❌ Invalid session or unauthorized".into())
                .unwrap());
        }
        Err(e) => {
            tracing::error!("❌ DB error checking session: {}", e);
            return Ok(Response::builder()
                .status(500)
                .header("content-type", "text/plain")
                .body("❌ Internal server error".into())
                .unwrap());
        }
    }

    // Start LLM stream
    let llm_stream = match query_ollama_stream(prompt.clone()).await {
        Ok(s) => s,
        Err(e) => {
            return Ok(Response::builder()
                .status(500)
                .header("content-type", "text/plain")
                .body(format!("❌ LLM Error: {e}\n").into())
                .unwrap());
        }
    };

    // Tee the stream into two sinks:
    //  - one for client
    //  - one to collect into Vec<Bytes> for saving
    let (client_tx, mut client_rx) = mpsc::channel::<Result<Bytes, std::io::Error>>(16);
    let (buffer_tx, mut buffer_rx) = mpsc::channel::<Bytes>(16);

    // Fork LLM stream and tee it
    tokio::spawn(async move {
        let mut stream = llm_stream;

        while let Some(chunk) = stream.next().await {
            if let Ok(ref bytes) = chunk {
                let _ = buffer_tx.send(bytes.clone()).await;
            }
            let _ = client_tx.send(chunk).await;
        }
    });

    // Stream to client
    let body_stream = stream! {
        while let Some(item) = client_rx.recv().await {
            yield item;
        }
    };

    let body = Body::from_stream(StreamBody::new(body_stream));
    let response = Response::builder()
        .status(200)
        .header("content-type", "text/event-stream")
        .body(body)
        .unwrap();

    // Store to DB after stream completes
    let pool = pool.clone();
    let prompt_clone = prompt.clone();

    tokio::spawn(async move {
        let mut collected = String::new();

        while let Some(chunk) = buffer_rx.recv().await {
            if let Ok(text) = std::str::from_utf8(&chunk) {
                if let Ok(json) = serde_json::from_str::<serde_json::Value>(text) {
                    if let Some(resp) = json.get("response").and_then(|r| r.as_str()) {
                        collected.push_str(resp);
                    }
                }
            }
        }

        if let Err(e) = sqlx::query_as::<_, Message>(
            r#"
            INSERT INTO messages (session_id, user_message, system_message)
            VALUES ($1, $2, $3)
            RETURNING id, created_at, session_id, user_message, system_message
            "#,
        )
        .bind(session_id)
        .bind(prompt_clone)
        .bind(collected)
        .fetch_one(&pool)
        .await
        {
            tracing::error!("❌ Failed to insert message: {}", e);
        }
    });

    Ok(response)
}
