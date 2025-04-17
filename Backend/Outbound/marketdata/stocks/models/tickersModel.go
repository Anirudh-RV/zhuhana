package stocks

type AllTickersAPIResponse struct {
	Count     int                         `json:"count"`
	NextURL   string                      `json:"next_url"`
	RequestID string                      `json:"request_id"`
	Results   []AllTickerStockInformation `json:"results"`
	Status    string                      `json:"status"`
}

type AllTickerStockInformation struct {
	Active          bool   `json:"active"`
	CIK             string `json:"cik"`
	CompositeFIGI   string `json:"composite_figi"`
	CurrencyName    string `json:"currency_name"`
	LastUpdatedUTC  string `json:"last_updated_utc"`
	Locale          string `json:"locale"`
	Market          string `json:"market"`
	Name            string `json:"name"`
	PrimaryExchange string `json:"primary_exchange"`
	ShareClassFIGI  string `json:"share_class_figi"`
	Ticker          string `json:"ticker"`
	Type            string `json:"type"`
}
