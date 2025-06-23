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
	StatusDescription string `json:"status_description"`
	OHLCData          []OHLC `json:"ohlc_data"`
	PageLimit         int    `json:"page_limit"`
	PageNo            int    `json:"page_no"`
	NextPageNo        *int   `json:"next_page_no,omitempty"`
}
