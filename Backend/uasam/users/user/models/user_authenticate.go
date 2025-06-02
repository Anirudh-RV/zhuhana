package models

type UserAuthenticateRequestHeaders struct {
	UserToken string `header:"USER_TOKEN" binding:"required"`
}

type UserAuthenticateResponse struct {
	Status            int         `json:"status"`
	StatusDescription string      `json:"statusDescription"`
	User              *UserObject `json:"user,omitempty"`
}
