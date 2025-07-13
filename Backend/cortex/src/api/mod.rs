use axum::{Router, routing::{get, post}};
use crate::state::AppState; // ✅ import your AppState
mod ask;
mod session;

use self::{ask::handle_ask, session::handle_create_session};

pub fn routes() -> Router<AppState> {
    Router::new()
        .route("/v1/ask/", get(handle_ask))
        .route("/v1/session/", post(handle_create_session))
}
