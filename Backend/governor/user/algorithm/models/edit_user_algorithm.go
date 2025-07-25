package models

type EditUserAlgorithmRequest struct {
	AlgorithmID   string `form:"algorithmID" binding:"required"`
	AlgorithmName string `form:"algorithmName" binding:"required"`
}

type EditUserAlgorithmResponse struct {
	Status            int            `json:"status"`
	StatusDescription string         `json:"statusDescription"`
	UserAlgorithm     *UserAlgorithm `json:"user_algorithm,omitempty"`
}
