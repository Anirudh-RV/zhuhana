package models

import "time"

type StartUserAlgorithmRequest struct {
	AlgorithmID   string     `json:"algorithmID" binding:"required"`
	Market        string     `json:"market" binding:"required"`
	Symbol        string     `json:"symbol" binding:"required"`
	StartTime     *time.Time `json:"start_time" binding:"required"`
	EndTime       *time.Time `json:"end_time" binding:"required"`
	Frequency     int        `json:"frequency" binding:"required"`
	PortfolioSize int        `json:"portfolio_size" binding:"required"`
}

type StartUserAlgorithmResponse struct {
	Status            int    `json:"status"`
	StatusDescription string `json:"statusDescription"`
}
