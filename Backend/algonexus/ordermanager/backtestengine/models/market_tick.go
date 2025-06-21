package models

import "time"

// MarketTick ticker,volume,open,close,high,low,window_start,transactions
type MarketTick struct {
	Ticker       string
	Volume       int
	Open         float64
	Close        float64
	High         float64
	Low          float64
	Start        time.Time
	Transactions float64
}
