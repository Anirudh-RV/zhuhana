package models

type MicroServiceAuthenticateRequestHeaders struct {
	AuthToken     string `header:"AUTH_TOKEN" binding:"required"`
	OriginService string `header:"ORIGIN_SERVICE" binding:"required"`
}

type MicroServiceAuthenticateResponse struct {
	Status            int    `json:"status"`
	StatusDescription string `json:"statusDescription"`
	CallerService     string `json:"callerService,omitempty"`
	CalleeService     string `json:"calleeService,omitempty"`
}
