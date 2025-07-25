package models

type ResetPasswordInitRequest struct {
	EmailID string `json:"emailId"`
}

type ResetPasswordInitResponse struct {
	Status            int    `json:"status"`
	StatusDescription string `json:"statusDescription"`
}

type ResetPasswordRequest struct {
	EmailID  string `json:"emailId"`
	Token    string `json:"token"`
	Password string `json:"password"`
}

type ResetPasswordResponse struct {
	Status            int    `json:"status"`
	StatusDescription string `json:"statusDescription"`
}

type UpdatePasswordRequest struct {
	Password string `json:"password"`
}

type UpdatePasswordResponse struct {
	Status            int    `json:"status"`
	StatusDescription string `json:"statusDescription"`
}
