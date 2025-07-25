package models

type UserSecretsDeleteRequest struct {
	ID string `json:"id" binding:"required"`
}

type UserSecretDeleteResponse struct {
	Status            int    `json:"status"`
	StatusDescription string `json:"statusDescription"`
}
