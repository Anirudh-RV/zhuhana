package models

import "time"

type OHLC struct {
	Symbol    string
	Market    string
	Date_Time time.Time
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    uint64
	Day       uint8
	Weekday   uint8
	Week      uint8
	Month     uint8
	Year      uint16
}
