use axum::{
    extract::{Request, State},
    http::StatusCode,
    middleware::Next,
    response::{IntoResponse, Response},
};
use redis::AsyncCommands;
use std::sync::Arc;
use crate::state::AppState;
use serde_json::json;

#[derive(Clone)]
pub struct RateLimiterConfig {
    pub source: String,       // "header", "query", or "body"
    pub param: String,        // parameter key, like "email"
    pub limit: i32,           // number of allowed requests
    pub window_secs: usize,   // duration of the window in seconds
    pub enable_param: bool,
    pub enable_ip: bool,
    pub ip_limit: i32,
    pub ip_window_secs: usize,
    pub endpoint: String,
}

pub async fn rate_limiter_middleware(
    State(state): State<AppState>,
    req: Request,
    next: Next,
    config: RateLimiterConfig,
) -> Result<Response, StatusCode> {
    let mut redis = state
    .redis
    .get_multiplexed_async_connection()
    .await
    .map_err(|_| StatusCode::INTERNAL_SERVER_ERROR)?;


    let ip = req
        .headers()
        .get("x-forwarded-for")
        .and_then(|v| v.to_str().ok())
        .unwrap_or("unknown_ip")
        .to_string();

    let mut pipe = redis::pipe();

    // IP rate limiting
    if config.enable_ip {
        let ip_key = format!("rl:{}:ip:{}", config.endpoint, ip);
        let ip_count: i32 = redis.get(&ip_key).await.unwrap_or(0);
        if ip_count >= config.ip_limit {
            return rate_limit_exceeded("IP rate limit exceeded");
        }
        pipe.incr(&ip_key, 1).ignore();
        if ip_count == 0 {
            pipe.expire(&ip_key, config.ip_window_secs.try_into().unwrap()).ignore();

        }
    }

    // Param-based rate limiting
    if config.enable_param {
        let param_val = match config.source.as_str() {
            "query" => req.uri().query().and_then(|q| {
                q.split('&')
                    .find(|s| s.starts_with(&config.param))
                    .and_then(|kv| kv.split('=').nth(1))
                    .map(|s| s.to_string())
            }),
            "header" => req
                .headers()
                .get(&config.param)
                .and_then(|v| v.to_str().ok())
                .map(|s| s.to_string()),
            _ => None,
        };

        if let Some(param_val) = param_val {
            let param_key = format!("rl:{}:param:{}:{}", config.endpoint, config.param, param_val);
            let param_count: i32 = redis.get(&param_key).await.unwrap_or(0);
            if param_count >= config.limit {
                return rate_limit_exceeded("User rate limit exceeded");
            }
            pipe.incr(&param_key, 1).ignore();
            if param_count == 0 {
                pipe.expire(&param_key, config.window_secs.try_into().unwrap()).ignore();

            }
        } else {
            return rate_limit_exceeded("Missing required parameter");
        }
    }

    let _ = pipe.query_async::<_, ()>(&mut redis).await;
    Ok(next.run(req).await)
}

fn rate_limit_exceeded(message: &str) -> Result<Response, StatusCode> {
    let body = axum::Json(json!({
        "status": -1,
        "statusDescription": message
    }));
    Ok((StatusCode::TOO_MANY_REQUESTS, body).into_response())
}
