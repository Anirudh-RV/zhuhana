use chrono::NaiveDateTime;
use uuid::Uuid;
use sqlx::FromRow;
use serde::{Deserialize, Serialize};


#[derive(Debug, Serialize, Deserialize, sqlx::FromRow)]
pub struct Session {
    pub id: Uuid,
    pub created_at: NaiveDateTime,
    pub user_id: Uuid,
    pub algorithm_id: Uuid,
    pub title: String,
}


#[derive(Debug, Serialize, Deserialize, sqlx::FromRow)]
pub struct Message {
    pub id: Uuid,
    pub session_id: Uuid,
    pub created_at: NaiveDateTime,
    pub user_message: String,
    pub system_message: String,
    pub model: String,
    pub tokens: i32,
}
