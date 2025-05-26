package models

type MicroServiceUserLoginRequestHeaders struct {
	APIKey        string `header:"API_KEY" binding:"required"`
	OriginService string `header:"ORIGIN_SERVICE" binding:"required"`
}

type MicroServiceUserLoginRequest struct {
	UserID string `header:"API_KEY" binding:"required"`
}

type MicroServiceUserLoginResponse struct {
	Status            int    `json:"status"`
	StatusDescription string `json:"statusDescription"`
	AccessToken       string `json:"accessToken,omitempty"`
}
