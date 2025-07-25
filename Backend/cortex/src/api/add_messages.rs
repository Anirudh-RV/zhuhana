use axum::{
    extract::{State, Extension, Json},
    http::StatusCode,
    response::IntoResponse,
};
use uuid::Uuid;
use chrono::Utc;

use crate::{
    api::{get_messages::ApiResponse, session::EmptyResult}, auth::types::UserObject, state::AppState
};
use serde::{Deserialize};


#[derive(Debug, Deserialize)]
pub struct PostMessageRequest {
    pub session_id: Uuid,
    pub system_q: String,
}


pub async fn handle_add_message(
    Extension(user): Extension<UserObject>,
    State(state): State<AppState>,
    Json(payload): Json<PostMessageRequest>,
) -> Result<impl IntoResponse, std::convert::Infallible> {
    let pool = &state.db;

    // Step 1: Validate session ownership
    let session_exists = sqlx::query_scalar::<_, Uuid>(
        r#"
        SELECT id FROM sessions
        WHERE id = $1 AND user_id = $2
        "#
    )
    .bind(payload.session_id)
    .bind(user.id)
    .fetch_optional(pool)
    .await;

    match session_exists {
        Ok(Some(_)) => {
            // Step 2: Insert the system message
            let message_id = Uuid::new_v4();
            let now = Utc::now().naive_utc();

            let insert_result = sqlx::query(
                r#"
                INSERT INTO messages (id, session_id, created_at, user_message, system_message, model, tokens)
                VALUES ($1, $2, $3, '', $4, '', 0)
                "#
            )
            .bind(message_id)
            .bind(payload.session_id)
            .bind(now)
            .bind(&payload.system_q)
            .execute(pool)
            .await;

            match insert_result {
                Ok(_) => {
                    let response = ApiResponse {
                        status: 1,
                        status_description: "System message added".to_string(),
                        result: EmptyResult {},
                    };
                    Ok((StatusCode::CREATED, Json(response)).into_response())
                }
                Err(e) => {
                    let response = ApiResponse {
                        status: 0,
                        status_description: format!("Failed to insert message: {e}"),
                        result: EmptyResult {},
                    };
                    Ok((StatusCode::INTERNAL_SERVER_ERROR, Json(response)).into_response())
                }
            }
        }
        Ok(None) => {
            let response = ApiResponse {
                status: 0,
                status_description: "Session not found or unauthorized".to_string(),
                result: EmptyResult {},
            };
            Ok((StatusCode::NOT_FOUND, Json(response)).into_response())
        }
        Err(e) => {
            let response = ApiResponse {
                status: 0,
                status_description: format!("Failed to validate session: {e}"),
                result: EmptyResult {},
            };
            Ok((StatusCode::INTERNAL_SERVER_ERROR, Json(response)).into_response())
        }
    }
}
