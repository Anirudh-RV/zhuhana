package models

type UserSecretsGetRequestHeaders struct {
	UserToken string `header:"USER_TOKEN" binding:"required"`
}

type UserSecretsGetRequest struct {
	Key string `form:"key" binding:"required"`
}

type UserSecretsGetResponse struct {
	Status            int         `json:"status"`
	StatusDescription string      `json:"statusDescription"`
	UserSecret        *UserSecret `json:"userSecret,omitempty"`
}
