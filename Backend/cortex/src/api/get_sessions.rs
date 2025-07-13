use axum::{
    extract::{Extension, State, Query},
    response::IntoResponse,
    http::StatusCode,
    Json,
};
use serde::{Serialize, Deserialize};
use uuid::Uuid;

use crate::{
    auth::types::UserObject,
    state::AppState,
    tables::Session,
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
pub struct GetSessionsParams {
    pub algorithm_id: Option<Uuid>,
}

pub async fn handle_get_sessions(
    Extension(user): Extension<UserObject>,
    State(state): State<AppState>,
    Query(params): Query<GetSessionsParams>,
) -> Result<impl IntoResponse, std::convert::Infallible> {
    let pool = &state.db;

    let result = if let Some(algorithm_id) = params.algorithm_id {
        sqlx::query_as::<_, Session>(
            r#"
            SELECT id, created_at, user_id, algorithm_id, title
            FROM sessions
            WHERE user_id = $1 AND algorithm_id = $2
            ORDER BY created_at DESC
            "#
        )
        .bind(user.id)
        .bind(algorithm_id)
        .fetch_all(pool)
        .await
    } else {
        sqlx::query_as::<_, Session>(
            r#"
            SELECT id, created_at, user_id, algorithm_id, title
            FROM sessions
            WHERE user_id = $1
            ORDER BY created_at DESC
            "#
        )
        .bind(user.id)
        .fetch_all(pool)
        .await
    };

    match result {
        Ok(sessions) => {
            let response = ApiResponse {
                status: 1,
                status_description: "Sessions fetched successfully".to_string(),
                result: sessions,
            };

            Ok((
                StatusCode::OK,
                [("content-type", "application/json")],
                Json(response),
            ).into_response())
        }
        Err(e) => {
            let response = ApiResponse {
                status: 0,
                status_description: format!("Failed to fetch sessions: {e}"),
                result: Vec::<Session>::new(),
            };

            Ok((
                StatusCode::INTERNAL_SERVER_ERROR,
                [("content-type", "application/json")],
                Json(response),
            ).into_response())
        }
    }
}
