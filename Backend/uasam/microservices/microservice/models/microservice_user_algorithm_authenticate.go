package models

type MicroServiceUserAlgorithmAuthenticateRequestHeaders struct {
	UserAlgorithmToken string `header:"USER_ALGORITHM_TOKEN" binding:"required"`
}

type MicroServiceUserAlgorithmAuthenticateResponse struct {
	Status            int    `json:"status"`
	StatusDescription string `json:"statusDescription"`
	UserAlgorithmID   string `json:"UserAlgorithmID,omitempty"`
}
