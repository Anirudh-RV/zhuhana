package models

type GetUserAlgorithmResponse struct {
	Status            int                `json:"status"`
	StatusDescription string             `json:"statusDescription"`
	UserAlgorithm     *UserAlgorithmInfo `json:"user_algorithm,omitempty"`
}

type GetAllUserAlgorithmsResponse struct {
	Status            int                 `json:"status"`
	StatusDescription string              `json:"statusDescription"`
	UserAlgorithms    []UserAlgorithmInfo `json:"user_algorithms"`
}
