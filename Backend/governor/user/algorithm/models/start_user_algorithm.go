package models

type StartUserAlgorithmCronScheduleRequest struct {
	AlgorithmID string `json:"algorithmID"`
}

type StartUserAlgorithmCronScheduleResponse struct {
	Status            int    `json:"status"`
	StatusDescription string `json:"statusDescription"`
}
