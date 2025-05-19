package commonutils

import (
	"encoding/json"
	"fmt"
	"net/http"
	constants "uasam/constants"
)

type IPGeoResponse struct {
	IPVersion     int     `json:"ipVersion"`
	IPAddress     string  `json:"ipAddress"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	CountryName   string  `json:"countryName"`
	CountryCode   string  `json:"countryCode"`
	TimeZone      string  `json:"timeZone"`
	ZipCode       string  `json:"zipCode"`
	CityName      string  `json:"cityName"`
	RegionName    string  `json:"regionName"`
	IsProxy       bool    `json:"isProxy"`
	Continent     string  `json:"continent"`
	ContinentCode string  `json:"continentCode"`
	Currency      struct {
		Code string `json:"code"`
		Name string `json:"name"`
	} `json:"currency"`
	Language  string   `json:"language"`
	TimeZones []string `json:"timeZones"`
	TLDs      []string `json:"tlds"`
}

func GetLocationForIpAddress(ipAddress string) (*IPGeoResponse, error) {
	url := fmt.Sprintf("%s%s", constants.IP_LOCATION_API_ENDPOINT, ipAddress)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Request error:", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Non-200 response:", resp.Status)
		return nil, err
	}

	var ipInfo IPGeoResponse
	if err := json.NewDecoder(resp.Body).Decode(&ipInfo); err != nil {
		fmt.Println("Decode error:", err)
		return nil, err
	}
	return &ipInfo, nil
}
