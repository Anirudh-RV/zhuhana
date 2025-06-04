package models

type UpdateUserAlgorithmCronScheduleRequest struct {
	AlgorithmID  string `json:"algorithmID"`
	CronSchedule string `json:"cronSchedule"`
}

type UpdateUserAlgorithmCronScheduleResponse struct {
	Status            int    `json:"status"`
	StatusDescription string `json:"statusDescription"`
}
