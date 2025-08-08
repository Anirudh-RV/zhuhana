package models

type GetAllUserAlgorithmRunsResponse struct {
	Status            int                `json:"status"`
	StatusDescription string             `json:"statusDescription"`
	UserAlgorithmRuns []UserAlgorithmRun `json:"user_algorithm_runs,omitempty"`
}
