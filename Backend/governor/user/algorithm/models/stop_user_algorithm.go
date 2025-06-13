package models

type StopUserAlgorithmRequest struct {
	AlgorithmID string `json:"algorithmID"`
}

type StopUserAlgorithmResponse struct {
	Status            int    `json:"status"`
	StatusDescription string `json:"statusDescription"`
}
