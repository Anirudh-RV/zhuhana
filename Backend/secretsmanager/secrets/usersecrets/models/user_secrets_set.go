package models

type UserSecretsSetRequestHeaders struct {
	UserToken string `header:"USER_TOKEN" binding:"required"`
}

type UserSecretsSetRequest struct {
	Key   string `header:"key" binding:"required"`
	Value string `header:"value" binding:"required"`
}

type UserSecretsSetResponse struct {
	Status            int    `json:"status"`
	StatusDescription string `json:"statusDescription"`
}
