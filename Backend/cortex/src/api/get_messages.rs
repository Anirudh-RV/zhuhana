use axum::{
    extract::{Extension, State, Query},
    response::IntoResponse,
    http::StatusCode,
    Json,
};
use serde::{Deserialize, Serialize};
use uuid::Uuid;

use crate::{
    auth::types::UserObject,
    state::AppState,
    tables::{Message},
};

#[derive(Serialize)]
pub struct ApiResponse<T> {
    #[serde(rename = "Status")]
    pub status: i32,

    #[serde(rename = "StatusDescription")]
    pub status_description: String,

    #[serde(rename = "Result")]
    pub result: T,
}

#[derive(Debug, Deserialize)]
pub struct GetMessagesParams {
    pub session_id: Option<Uuid>,
}

pub async fn handle_get_messages(
    Extension(user): Extension<UserObject>,
    State(state): State<AppState>,
    Query(params): Query<GetMessagesParams>,
) -> Result<impl IntoResponse, std::convert::Infallible> {
    let pool = &state.db;

    let result = if let Some(session_id) = params.session_id {
        // First, validate session ownership
        let session_exists = sqlx::query_scalar::<_, Uuid>(
            r#"
            SELECT id FROM sessions
            WHERE id = $1 AND user_id = $2
            "#
        )
        .bind(session_id)
        .bind(user.id)
        .fetch_optional(pool)
        .await;

        match session_exists {
            Ok(Some(_)) => {
                // Fetch messages for this session
                sqlx::query_as::<_, Message>(
                    r#"
                    SELECT id, session_id, created_at, user_message, system_message, model, tokens
                    FROM messages
                    WHERE session_id = $1
                    ORDER BY created_at ASC
                    "#
                )
                .bind(session_id)
                .fetch_all(pool)
                .await
            }
            Ok(None) => {
                // Session not found or not owned by user
                return Ok((StatusCode::NOT_FOUND, Json(ApiResponse {
                    status: 0,
                    status_description: "Session not found or not authorized".to_string(),
                    result: Vec::<Message>::new(),
                })).into_response());
            }
            Err(e) => {
                return Ok((StatusCode::INTERNAL_SERVER_ERROR, Json(ApiResponse {
                    status: 0,
                    status_description: format!("DB error validating session: {e}"),
                    result: Vec::<Message>::new(),
                })).into_response());
            }
        }
    } else {
        // session_id is required; return 400
        return Ok((StatusCode::BAD_REQUEST, Json(ApiResponse {
            status: 0,
            status_description: "Missing required parameter: session_id".to_string(),
            result: Vec::<Message>::new(),
        })).into_response());
    };

    match result {
        Ok(messages) => {
            let response = ApiResponse {
                status: 1,
                status_description: "Messages fetched successfully".to_string(),
                result: messages,
            };

            Ok((StatusCode::OK, Json(response)).into_response())
        }
        Err(e) => {
            let response = ApiResponse {
                status: 0,
                status_description: format!("Failed to fetch messages: {e}"),
                result: Vec::<Message>::new(),
            };

            Ok((StatusCode::INTERNAL_SERVER_ERROR, Json(response)).into_response())
        }
    }
}
