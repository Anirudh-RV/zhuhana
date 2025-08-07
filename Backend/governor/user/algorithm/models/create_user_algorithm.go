package models

type CreateUserAlgorithmRequest struct {
	AlgorithmName string `form:"algorithmName" binding:"required"`
}

type CreateUserAlgorithmResponse struct {
	Status            int            `json:"status"`
	StatusDescription string         `json:"statusDescription"`
	UserAlgorithm     *UserAlgorithm `json:"user_algorithm,omitempty"`
}

type PythonBuilderRequest struct {
	UserID    string `json:"userID"`
	ScriptID  string `json:"scriptID"`
	ScriptURL string `json:"scriptURL"`
}

type PythonBuilderResponse struct {
	Status            int    `json:"status"`
	StatusDescription string `json:"statusDescription"`
}
