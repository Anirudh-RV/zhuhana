package polygon

import (
	"outbound/logger"
	polygonTickersService "outbound/marketdata/stocks/polygon/services"

	"github.com/gin-gonic/gin"
)

type PolygonTickersController struct {
	polygonTickersService *polygonTickersService.PolygonTickersService
	log                   *logger.Logger
}

func NewPolygonTickersController(polygonTickersService *polygonTickersService.PolygonTickersService, log *logger.Logger) *PolygonTickersController {
	return &PolygonTickersController{
		polygonTickersService: polygonTickersService,
		log:                   log,
	}
}

// @Summary Get all tickers
// @Description Fetches all tickers with an optional limit parameter
// @Tags Tickers
// @Accept json
// @Produce json
// @Param limit query int false "Number of tickers to retrieve (default: 10)"
// @Success 200 {object} tickerModels.<MAKE THIS> "List of tickers"
// @Failure 400 {object} map[string]string "Error message"
// @Router /api/marketdata/v1/stocks/polygon/ticker/ [get]
func (ptc *PolygonTickersController) GetDailyTickerOHLCV_V1(c *gin.Context) {
	// IMPLEMENT: https://polygon.io/docs/rest/stocks/aggregates/daily-ticker-summary
}
