package models

type SignUpInitRequest struct {
	FirstName  string  `json:"firstName"`
	MiddleName *string `json:"middleName,omitempty"`
	LastName   string  `json:"lastName"`
	EmailID    string  `json:"emailId"`
	Password   string  `json:"password"`
}

type SignUpInitResponse struct {
	Status            int    `json:"status"`
	StatusDescription string `json:"statusDescription"`
}

type SignUpVerifyOTPRequest struct {
	EmailID string `json:"emailId"`
	Otp     string `json:"otp"`
}

type SignUpVerifyOTPResponse struct {
	Status            int         `json:"status"`
	StatusDescription string      `json:"statusDescription"`
	User              *UserObject `json:"user,omitempty"`
	AccessToken       string      `json:"accessToken,omitempty"`
}
