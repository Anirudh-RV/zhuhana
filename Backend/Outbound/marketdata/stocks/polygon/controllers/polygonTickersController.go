package polygon

import (
	"net/http"
	"outbound/logger"
	tickerModels "outbound/marketdata/stocks/models"
	polygonTickersService "outbound/marketdata/stocks/polygon/services"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
// @Success 200 {object} tickerModels.AllTickersAPIResponse "List of tickers"
// @Failure 400 {object} map[string]string "Error message"
// @Router /api/marketdata/v1/stocks/polygon/all-tickers [get]
func (ptc *PolygonTickersController) GetAllTickersV1(c *gin.Context) {
	var tickers *tickerModels.AllTickersAPIResponse
	var err error
	param := c.DefaultQuery("limit", "10")
	limit, _ := strconv.Atoi(param)
	ptc.log.Info("AllTickersV1 api called with:", zap.String("execution level", "controller"), zap.String("limit", strconv.Itoa(limit)))

	tickers, err = ptc.polygonTickersService.GetAllTickersV1(limit)

	if err != nil {
		ptc.log.Error("error", zap.String("execution level", "controller"), zap.String("error", err.Error()))
		c.JSON(http.StatusBadRequest, err)
	}
	if tickers != nil {
		c.JSON(http.StatusOK, tickers)
	} else {
		ptc.log.Error("Empty Tickers", zap.String("execution level", "controller"))
		c.JSON(http.StatusBadRequest, nil)
	}
}
