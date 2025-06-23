package models

type MicroServiceUserAlgorithmLoginRequestHeaders struct {
	APIKey        string `header:"API_KEY" binding:"required"`
	OriginService string `header:"ORIGIN_SERVICE" binding:"required"`
}

type MicroServiceUserAlgorithmLoginRequest struct {
	UserAlgorithmID string `header:"API_KEY" binding:"required"`
}

type MicroServiceUserAlgorithmLoginResponse struct {
	Status            int    `json:"status"`
	StatusDescription string `json:"statusDescription"`
	AccessToken       string `json:"accessToken,omitempty"`
}
