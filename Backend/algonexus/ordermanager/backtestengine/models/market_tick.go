package models

import "time"

// MarketTick ticker,volume,open,close,high,low,window_start,transactions
type MarketTick struct {
	Ticker       string    `json:"ticker"`
	Volume       int       `json:"volume"`
	Open         float64   `json:"open"`
	Close        float64   `json:"close"`
	High         float64   `json:"high"`
	Low          float64   `json:"low"`
	Start        time.Time `json:"start"`
	Transactions int       `json:"transactions"`
}
