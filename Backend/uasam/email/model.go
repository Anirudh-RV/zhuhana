package email

type EmailAddress struct {
	Email string `json:"email"`
}

type SenderInfo struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type EmailRequest struct {
	Sender      SenderInfo     `json:"sender"`
	To          []EmailAddress `json:"to"`
	Subject     string         `json:"subject"`
	HTMLContent string         `json:"htmlContent"`
}
