use axum::{
    http::{Request, StatusCode},
    middleware::Next,
    response::Response,
    Extension,
};
use axum::body::Body;
use std::sync::Arc;
use reqwest::Client;

use crate::auth::types::{UserAuthenticateResponse, UserObject, ErrorResponse};

#[derive(Clone)]
pub struct AuthConfig {
    pub auth_service_url: String,
    pub http_client: Client,
}

fn unauthorized_response(message: &str) -> Response {
    // Create an error response object
    let err = ErrorResponse {
        status: -1,
        statusDescription: message.to_string(),
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
        statusDescription: message.to_string(),
    };
    let body = serde_json::to_string(&err).unwrap();
    Response::builder()
        .status(StatusCode::INTERNAL_SERVER_ERROR)
        .header("content-type", "application/json")
        .body(Body::from(body))
        .unwrap()
}

pub async fn user_auth_middleware(
    Extension(config): Extension<Arc<AuthConfig>>,
    mut req: Request<Body>,
    next: Next,
) -> Result<Response, StatusCode> {
    let user_token = match req.headers().get("USER_TOKEN") {
        Some(t) => t.to_str().unwrap_or("").to_string(),
        None => return Ok(unauthorized_response("Missing USER_TOKEN")),
    };

    let client = &config.http_client;

    let resp = client
    .post(&config.auth_service_url)
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

    let auth_response: UserAuthenticateResponse = match resp.json().await {
        Ok(body) => body,
        Err(err) => {
            tracing::error!("❌ Failed to parse auth response: {:?}", err);
            return Err(StatusCode::INTERNAL_SERVER_ERROR);
        }
    };


    if auth_response.status != 1 {
        return Ok(unauthorized_response(&auth_response.statusDescription));
    }

    if let Some(user) = auth_response.user {
        req.extensions_mut().insert(user);
    }

    Ok(next.run(req).await)
}
