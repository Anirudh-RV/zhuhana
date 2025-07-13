use serde::{Deserialize, Serialize};
use uuid::Uuid;
use chrono::{DateTime, Utc};


#[derive(Debug, Deserialize, Clone)]
pub struct UserObject {
    #[serde(rename = "ID")]
    pub id: Uuid,
    #[serde(rename = "FirstName")]
    pub first_name: String,
    #[serde(rename = "MiddleName")]
    pub middle_name: String,
    #[serde(rename = "LastName")]
    pub last_name: String,
    #[serde(rename = "EmailID")]
    pub email_id: String,
    #[serde(rename = "CreatedAt")]
    pub created_at: DateTime<Utc>,
    #[serde(rename = "UpdatedAt")]
    pub updated_at: DateTime<Utc>,
}

#[derive(Debug, Deserialize)]
pub struct UserAuthenticateResponse {
    pub status: i32,
    #[serde(rename = "statusDescription")]
    pub status_description: String,
    pub user: Option<UserObject>,
}


#[derive(Debug, Serialize)]
pub struct ErrorResponse {
    pub status: i32,
    #[serde(rename = "statusDescription")]
    pub status_description: String,
}
