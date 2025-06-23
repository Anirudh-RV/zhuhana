package models

import "time"

type OHLCRangeRequest struct {
	Symbol    string    `json:"symbol"`
	Market    string    `json:"market"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

type OHLCRangeResponse struct {
	Status            int    `json:"status"`
	StatusDescription string `json:"statusDescription"`
	OHLCData          []OHLC `json:"OHLCData,omitempty"`
}
