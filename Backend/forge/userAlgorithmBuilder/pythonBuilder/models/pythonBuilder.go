package models

type PythonBuilderRequest struct {
	UserID    string `json:"userID"`
	ScriptID  string `json:"scriptID"`
	ScriptURL string `json:"scriptURL"`
}

type PythonBuilderResponse struct {
	Status            int    `json:"status"`
	StatusDescription string `json:"statusDescription"`
}
