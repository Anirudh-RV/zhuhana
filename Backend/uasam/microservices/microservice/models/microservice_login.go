package models

type MicroServiceLoginRequestHeaders struct {
	APIKey        string `header:"API_KEY" binding:"required"`
	OriginService string `header:"ORIGIN_SERVICE" binding:"required"`
}

type MicroServiceLoginResponse struct {
	Status            int    `json:"status"`
	StatusDescription string `json:"statusDescription"`
	AccessToken       string `json:"accessToken,omitempty"`
}
