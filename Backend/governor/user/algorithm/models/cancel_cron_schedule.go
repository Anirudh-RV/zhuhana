package models

type CancelUserAlgorithmCronScheduleRequest struct {
	AlgorithmID string `json:"algorithmID"`
}

type CancelUserAlgorithmCronScheduleResponse struct {
	Status            int    `json:"status"`
	StatusDescription string `json:"statusDescription"`
}
