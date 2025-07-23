use axum::{
    extract::{Request, State},
    http::StatusCode,
    middleware::Next,
    response::Response,
    body::Body,
};
use reqwest::Client;

use crate::auth::types::{UserAuthenticateResponse, ErrorResponse};
use crate::{state::AppState};

#[derive(Clone)]
pub struct AuthConfig {
    pub auth_service_url: String,
    pub http_client: Client,
}

fn unauthorized_response(message: &str) -> Response {
    // Create an error response object
    let err = ErrorResponse {
        status: -1,
        status_description: message.to_string(),
    };
    let body = serde_json::to_string(&err).unwrap();
    Response::builder()
        .status(StatusCode::UNAUTHORIZED)
        .header("content-type", "application/json")
        .body(Body::from(body))
        .unwrap()
}

fn internal_error_response(message: &str) -> Response {
    let err = ErrorResponse {
        status: 0,
        status_description: message.to_string(),
    };
    let body = serde_json::to_string(&err).unwrap();
    Response::builder()
        .status(StatusCode::INTERNAL_SERVER_ERROR)
        .header("content-type", "application/json")
        .body(Body::from(body))
        .unwrap()
}

pub async fn user_auth_middleware(
    State(state): State<AppState>,
    mut req: Request<Body>,
    next: Next,
) -> Result<Response, StatusCode> {
    let user_token = match req.headers().get("USER_TOKEN") {
        Some(t) => t.to_str().unwrap_or("").to_string(),
        None => return Ok(unauthorized_response("Missing USER_TOKEN")),
    };

    let client = &state.auth_config.http_client;
    let auth_url = &state.auth_config.auth_service_url;

    let resp = client
        .post(auth_url)
        .header("USER_TOKEN", user_token)
        .send()
        .await
        .map_err(|err| {
            tracing::error!("❌ HTTP request to auth service failed: {:?}", err);
            StatusCode::UNAUTHORIZED
        })?;

    if !resp.status().is_success() {
        return Ok(unauthorized_response("Not Authorized"));
    }

    let raw_body = match resp.text().await {
        Ok(body) => {
            tracing::error!("🔍 Raw auth response: {}", body);
            body
        }
        Err(err) => {
            tracing::error!("❌ Failed to read auth response body: {:?}", err);
            return Err(StatusCode::INTERNAL_SERVER_ERROR);
        }
    };

    let auth_response: UserAuthenticateResponse = match serde_json::from_str(&raw_body) {
        Ok(body) => body,
        Err(err) => {
            tracing::error!("❌ Failed to parse auth response: {:?}", err);
            tracing::error!("📦 Raw response that failed to parse: {}", raw_body);
            return Err(StatusCode::INTERNAL_SERVER_ERROR);
        }
    };


    if auth_response.status != 1 {
        return Ok(unauthorized_response(&auth_response.status_description));
    }

    if let Some(user) = auth_response.user {
        req.extensions_mut().insert(user);
    }

    Ok(next.run(req).await)
}
