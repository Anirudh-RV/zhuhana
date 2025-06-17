package models

type StartUserAlgorithmRequest struct {
	AlgorithmID string `json:"algorithmID"`
}

type StartUserAlgorithmResponse struct {
	Status            int    `json:"status"`
	StatusDescription string `json:"statusDescription"`
}
