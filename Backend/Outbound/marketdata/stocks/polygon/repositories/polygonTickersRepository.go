package polygon

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"outbound/logger"
	tickerModels "outbound/marketdata/stocks/models"
	polygonConstants "outbound/marketdata/stocks/polygon/constants"
	"strconv"

	"go.uber.org/zap"
)

type PolygonTickersRepository struct {
	log *logger.Logger
}

func NewPolygonTickersRepository(log *logger.Logger) *PolygonTickersRepository {
	return &PolygonTickersRepository{
		log: log,
	}
}

type TickersRepositoryInterface interface {
	GetAllTickersV1(limit int) (*tickerModels.AllTickersAPIResponse, error)
}

func (ptr *PolygonTickersRepository) GetAllTickersV1(limit int) (*tickerModels.AllTickersAPIResponse, error) {
	// TODO: Handle errors in the response
	params := url.Values{}
	params.Add("market", "stocks")
	params.Add("active", "true")
	params.Add("order", "asc")
	params.Add("limit", strconv.Itoa(limit))
	params.Add("sort", "ticker")
	params.Add("apiKey", "USE_FROM_HEADER")

	finalURL := fmt.Sprintf("%s?%s", polygonConstants.AllTickersBaseURL, params.Encode())
	ptr.log.Info("AllTickersV1 API:", zap.String("execution level", "repository"), zap.String("url", finalURL))

	// Create a GET request
	resp, err := http.Get(finalURL)
	if err != nil {
		ptr.log.Error("error making request", zap.String("execution level", "repository"), zap.String("error", err.Error()))
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ptr.log.Error("error reading response", zap.String("execution level", "repository"), zap.String("error", err.Error()))
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	// Unmarshal JSON response into struct
	var apiResponse tickerModels.AllTickersAPIResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		ptr.log.Error("error unmarshaling json", zap.String("execution level", "Repository"), zap.String("error", err.Error()))
		return nil, fmt.Errorf("error unmarshaling json: %v", err)
	}

	return &apiResponse, nil
}
