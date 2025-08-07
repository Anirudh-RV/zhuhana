use axum::{
    extract::{State, Json, Extension},
    http::StatusCode,
    response::{IntoResponse},
    Json as AxumJson,
};
use uuid::Uuid;
use serde::{Serialize, Deserialize};

use crate::{
    auth::types::UserObject,
    tables::Session,
    state::AppState,
};

#[derive(Debug, Deserialize)]
pub struct CreateSessionRequest {
    pub algorithm_id: Option<Uuid>,
    pub title: String,
    pub system_q: Option<String>,
}

#[derive(Serialize)]
pub struct ApiResponse<T> {
    #[serde(rename = "Status")]
    pub status: i32,

    #[serde(rename = "StatusDescription")]
    pub status_description: String,

    #[serde(rename = "Result")]
    pub result: T,
}

#[derive(Serialize)]
pub struct EmptyResult {}

pub async fn handle_create_session(
    Extension(user): Extension<UserObject>,
    State(state): State<AppState>,
    Json(payload): Json<CreateSessionRequest>,
) -> Result<impl IntoResponse, std::convert::Infallible> {
    let pool = &state.db;

    let maybe_session = sqlx::query_as::<_, Session>(
        r#"
        INSERT INTO sessions (user_id, algorithm_id, title)
        VALUES ($1, $2, $3)
        RETURNING id, created_at, user_id, algorithm_id, title
        "#
    )
    .bind(user.id)
    .bind(payload.algorithm_id)
    .bind(&payload.title)
    .fetch_one(pool)
    .await;

    match maybe_session {
        Ok(session) => {
            if let Some(system_q) = payload.system_q {
                // Insert the initial system message
                let _ = sqlx::query(
                    r#"
                    INSERT INTO messages (id, session_id, created_at, user_message, system_message, model, tokens)
                    VALUES ($1, $2, NOW(), '', $3, '', 0)
                    "#
                )
                .bind(Uuid::new_v4())
                .bind(session.id)
                .bind(system_q)
                .execute(pool)
                .await;
            }

            let response = ApiResponse {
                status: 1,
                status_description: "Session Created".to_string(),
                result: session,
            };

            Ok((
                StatusCode::CREATED,
                [("content-type", "application/json")],
                AxumJson(response),
            ).into_response())
        }

        Err(e) => {
            let response = ApiResponse {
                status: 0,
                status_description: format!("Failed to create session: {e}"),
                result: EmptyResult {},
            };

            Ok((
                StatusCode::INTERNAL_SERVER_ERROR,
                [("content-type", "application/json")],
                AxumJson(response),
            ).into_response())
        }
    }
}
