package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"uasam/news/constants"
	"uasam/news/models"

	"uasam/logger"

	"go.uber.org/zap"
)

type NewsArticleService struct {
	logger                *logger.Logger
	newsAPIKEY            string
	getNewsArticleBaseURL string
	articlePageLimit      int
	countries             string
}

func NewNewsArticleService(logger *logger.Logger) *NewsArticleService {
	newsAPIKey := os.Getenv("NEWS_API_KEY")
	newsAPIEndpoint := os.Getenv("NEWS_API_ENDPOINT")
	return &NewsArticleService{
		logger:                logger,
		newsAPIKEY:            newsAPIKey,
		getNewsArticleBaseURL: newsAPIEndpoint,
		articlePageLimit:      constants.GET_NEWS_ARTICLES_LIMIT,
		countries:             "us,sg,in,cn",
	}
}

func (nas *NewsArticleService) joinCountries(countries []string) string {
	return strings.Join(countries, ",")
}

func (nas *NewsArticleService) GetNewsArticle(query string, nextPage string) (*models.NewsArticleData, error) {
	// Build query params
	params := url.Values{}
	params.Set("apikey", nas.newsAPIKEY)
	params.Set("q", query)
	params.Set("size", fmt.Sprintf("%d", nas.articlePageLimit))
	if nas.countries != "" {
		params.Set("country", nas.countries)
	}
	if nextPage != "" {
		params.Set("page", nextPage)
	}

	fullURL := fmt.Sprintf("%s?%s", nas.getNewsArticleBaseURL, params.Encode())

	// Send GET request
	resp, err := http.Get(fullURL)
	if err != nil {
		go nas.logger.Warning("http request error", zap.String("execution level", "GetNewsArticle"), zap.String("Error", err.Error()))
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Read body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		go nas.logger.Warning("failed to read response", zap.String("execution level", "GetNewsArticle"), zap.String("Error", err.Error()))
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	nas.logger.Info("Raw response", zap.String("body", string(body)))

	// Unmarshal
	var result models.NewsDataResponse
	if err := json.Unmarshal(body, &result); err != nil {
		go nas.logger.Warning("failed to unmarshal json", zap.String("execution level", "GetNewsArticle"), zap.String("Error", err.Error()))
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return &models.NewsArticleData{
		Results:  result.Results,
		NextPage: result.NextPage,
	}, nil
}
