package models

type MicroServiceUserAuthenticateRequestHeaders struct {
	UserScriptToken string `header:"USER_SERVICE_TOKEN" binding:"required"`
}

type MicroServiceUserAuthenticateResponse struct {
	Status            int    `json:"status"`
	StatusDescription string `json:"statusDescription"`
	UserID            string `json:"UserID,omitempty"`
}
