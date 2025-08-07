package models

type UpdateUserFieldsRequest struct {
	FirstName  *string `json:"first_name"`
	MiddleName *string `json:"middle_name"`
	LastName   *string `json:"last_name"`
}

type UpdateUserFieldsResponse struct {
	Status            int    `json:"status"`
	StatusDescription string `json:"statusDescription"`
}
