package stocks

import (
	"encoding/json"
	"fmt"
	"io"
	"marketdata/constants"
	"marketdata/logger"
	stocksConstants "marketdata/stocks/constants"
	tickerModels "marketdata/stocks/models"
	"net/http"
	"net/url"
	"strconv"

	"go.uber.org/zap"
)

type TickersRepository struct {
	log *logger.Logger
}

func NewTickersRepository(log *logger.Logger) *TickersRepository {
	return &TickersRepository{
		log: log,
	}
}

type TickersRepositoryInterface interface {
	GetAllTickersV1(limit int) (*tickerModels.AllTickersAPIResponse, error)
}

func (tickersRepo *TickersRepository) GetAllTickersV1(limit int) (*tickerModels.AllTickersAPIResponse, error) {
	// TODO: Handle errors in the response
	params := url.Values{}
	params.Add("market", "stocks")
	params.Add("active", "true")
	params.Add("order", "asc")
	params.Add("limit", strconv.Itoa(limit))
	params.Add("sort", "ticker")
	params.Add("apiKey", constants.POLYGON_API_KEY)

	finalURL := fmt.Sprintf("%s?%s", stocksConstants.AllTickersBaseURL, params.Encode())
	tickersRepo.log.Info("AllTickersV1 API:", zap.String("execution level", "repository"), zap.String("url", finalURL))

	// Create a GET request
	resp, err := http.Get(finalURL)
	if err != nil {
		tickersRepo.log.Error("error making request", zap.String("execution level", "repository"), zap.String("error", err.Error()))
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		tickersRepo.log.Error("error reading response", zap.String("execution level", "repository"), zap.String("error", err.Error()))
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	// Unmarshal JSON response into struct
	var apiResponse tickerModels.AllTickersAPIResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		tickersRepo.log.Error("error unmarshaling json", zap.String("execution level", "Repository"), zap.String("error", err.Error()))
		return nil, fmt.Errorf("error unmarshaling json: %v", err)
	}

	return &apiResponse, nil
}
