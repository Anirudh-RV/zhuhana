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
	StatusDescription string `json:"status_description"`
}
