package models

type UserSecretsGetRequestHeaders struct {
	UserScriptToken string `header:"USER_SCRIPT_TOKEN" binding:"required"`
}

type UserSecretsGetRequest struct {
	Key string `form:"key" binding:"required"`
}

type UserSecretsGetResponse struct {
	Status            int         `json:"status"`
	StatusDescription string      `json:"statusDescription"`
	UserSecret        *UserSecret `json:"userSecret,omitempty"`
}
