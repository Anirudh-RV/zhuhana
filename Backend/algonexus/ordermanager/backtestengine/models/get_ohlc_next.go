package models

import "time"

type OHLCNextRequest struct {
	Symbol      string    `json:"symbol"`
	Market      string    `json:"market"`
	CurrentTime time.Time `json:"current_time"`
	EndTime     time.Time `json:"end_time"`
	NextStep    int       `json:"next_step"`
}

type OHLCNextResponse struct {
	Status            int     `json:"status"`
	StatusDescription string  `json:"status_description"`
	OHLCData          *OHLC   `json:"ohlc_data,omitempty"`
	NextURL           *string `json:"next_url,omitempty"`
}
