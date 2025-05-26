package models

type MicroServiceUserAuthenticateRequestHeaders struct {
	UserToken string `header:"USER_TOKEN" binding:"required"`
}

type MicroServiceUserAuthenticateResponse struct {
	Status            int    `json:"status"`
	StatusDescription string `json:"statusDescription"`
	UserID            string `json:"UserID,omitempty"`
}
