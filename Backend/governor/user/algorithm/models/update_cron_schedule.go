package models

type UpdateUserAlgorithmCronScheduleRequest struct {
	AlgorithmID       string `json:"algorithmID"`
	StartCronSchedule string `json:"startCronSchedule"`
	EndCronSchedule   string `json:"endCronSchedule"`
}

type UpdateUserAlgorithmCronScheduleResponse struct {
	Status            int    `json:"status"`
	StatusDescription string `json:"statusDescription"`
}
