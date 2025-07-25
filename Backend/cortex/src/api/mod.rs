use axum::{Router, routing::{get, post}};
use crate::{api::{get_messages::handle_get_messages, get_sessions::handle_get_sessions, add_messages::handle_add_message}, state::AppState}; // ✅ import your AppState
mod ask;
mod session;
mod get_sessions;
mod get_messages;
mod add_messages;

use self::{ask::handle_ask, session::handle_create_session};

pub fn routes() -> Router<AppState> {
    Router::new()
        .route("/v1/ask/", get(handle_ask))
        .route("/v1/session/", post(handle_create_session))
        .route("/v1/session/", get(handle_get_sessions))
        .route("/v1/messages/", get(handle_get_messages))
        .route("/v1/message/", post(handle_add_message))
}
