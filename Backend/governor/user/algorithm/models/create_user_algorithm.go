package models

type CreateUserAlgorithmRequest struct {
	ScriptName   string `form:"scriptName" binding:"required"`
	CronSchedule string `form:"cronSchedule" binding:"required"`
}

type CreateUserAlgorithmResponse struct {
	Status            int            `json:"status"`
	StatusDescription string         `json:"statusDescription"`
	UserAlgorithm     *UserAlgorithm `json:"user_algorithm,omitempty"`
}
