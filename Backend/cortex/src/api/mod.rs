use axum::{
    middleware,
    routing::{get, post},
    Router,
};
use crate::{api::{get_messages::handle_get_messages, get_sessions::handle_get_sessions, add_messages::handle_add_message, rate_limiter::{rate_limiter_middleware, RateLimiterConfig}}, state::AppState}; // ✅ import your AppState
mod ask;
mod session;
mod get_sessions;
mod get_messages;
mod add_messages;
mod rate_limiter;

use self::{ask::handle_ask, session::handle_create_session};

pub fn routes(state: AppState) -> Router<AppState> {
    let ask_route = Router::new().route(
        "/v1/ask/",
        get(handle_ask).route_layer(middleware::from_fn_with_state(
            state.clone(),
            |state, req, next| {
                rate_limiter_middleware(
                    state,
                    req,
                    next,
                    RateLimiterConfig {
                        source: "header".into(),
                        param: "USER_TOKEN".into(),
                        limit: 100,
                        window_secs: 300,
                        enable_param: true,
                        enable_ip: true,
                        ip_limit: 300,
                        ip_window_secs: 300,
                        endpoint: "ask".into(),
                    },
                )
            },
        )),
    );

    let session_post = Router::new().route(
        "/v1/session/",
        post(handle_create_session).route_layer(middleware::from_fn_with_state(
            state.clone(),
            |state, req, next| {
                rate_limiter_middleware(
                    state,
                    req,
                    next,
                    RateLimiterConfig {
                        source: "header".into(),
                        param: "USER_TOKEN".into(),
                        limit: 100,
                        window_secs: 300,
                        enable_param: true,
                        enable_ip: true,
                        ip_limit: 300,
                        ip_window_secs: 300,
                        endpoint: "create_session".into(),
                    },
                )
            },
        )),
    );

    let session_get = Router::new().route(
        "/v1/session/",
        get(handle_get_sessions).route_layer(middleware::from_fn_with_state(
            state.clone(),
            |state, req, next| {
                rate_limiter_middleware(
                    state,
                    req,
                    next,
                    RateLimiterConfig {
                        source: "header".into(),
                        param: "USER_TOKEN".into(),
                        limit: 100,
                        window_secs: 300,
                        enable_param: true,
                        enable_ip: true,
                        ip_limit: 300,
                        ip_window_secs: 300,
                        endpoint: "get_sessions".into(),
                    },
                )
            },
        )),
    );

    let get_messages = Router::new().route(
        "/v1/messages/",
        get(handle_get_messages).route_layer(middleware::from_fn_with_state(
            state.clone(),
            |state, req, next| {
                rate_limiter_middleware(
                    state,
                    req,
                    next,
                    RateLimiterConfig {
                        source: "header".into(),
                        param: "USER_TOKEN".into(),
                        limit: 300,
                        window_secs: 300,
                        enable_param: true,
                        enable_ip: true,
                        ip_limit: 300,
                        ip_window_secs: 300,
                        endpoint: "get_messages".into(),
                    },
                )
            },
        )),
    );

    let add_message = Router::new().route(
        "/v1/message/",
        post(handle_add_message).route_layer(middleware::from_fn_with_state(
            state.clone(),
            |state, req, next| {
                rate_limiter_middleware(
                    state,
                    req,
                    next,
                    RateLimiterConfig {
                        source: "header".into(),
                        param: "USER_TOKEN".into(),
                        limit: 100,
                        window_secs: 300,
                        enable_param: true,
                        enable_ip: true,
                        ip_limit: 300,
                        ip_window_secs: 300,
                        endpoint: "add_message".into(),
                    },
                )
            },
        )),
    );

    ask_route
        .merge(session_post)
        .merge(session_get)
        .merge(get_messages)
        .merge(add_message)
        .with_state(state)
}
