package models

type LoginVerifyPasswordRequest struct {
	EmailID  string `json:"emailId"`
	Password string `json:"password"`
}

type LoginVerifyPasswordResponse struct {
	Status            int    `json:"status"`
	StatusDescription string `json:"statusDescription"`
}

type LoginVerifyOTPRequest struct {
	EmailID string `json:"emailId"`
	Otp     string `json:"otp"`
}

type LoginVerifyOTPResponse struct {
	Status            int         `json:"status"`
	StatusDescription string      `json:"statusDescription"`
	User              *UserObject `json:"user,omitempty"`
	AccessToken       string      `json:"accessToken,omitempty"`
}
