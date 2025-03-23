package stocks

import (
	"fmt"
	"marketdata/logger"
	tickerModels "marketdata/stocks/models"
	tickersService "marketdata/stocks/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type TickersController struct {
	tickersService *tickersService.TickersService
	log            *logger.Logger
}

func NewTickersController(tickersService *tickersService.TickersService, log *logger.Logger) *TickersController {
	return &TickersController{
		tickersService: tickersService,
		log:            log,
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
// @Router /api/marketdata/v1/stocks/all-tickers [get]
func (tc *TickersController) GetAllTickersV1(c *gin.Context) {
	var tickers *tickerModels.AllTickersAPIResponse
	var err error
	param := c.DefaultQuery("limit", "10")
	limit, _ := strconv.Atoi(param)
	tc.log.Info("AllTickersV1 API called with:", zap.String("Execution Level", "Controller"), zap.String("limit", strconv.Itoa(limit)))

	tickers, err = tc.tickersService.GetAllTickersV1(limit)

	if err != nil {
		fmt.Printf("ERROR - %s", err)
		tc.log.Error("Error", zap.String("Execution Level", "Controller"), zap.String("Error", err.Error()))
		c.JSON(http.StatusBadRequest, err)
	}
	if tickers != nil {
		c.JSON(http.StatusOK, tickers)
	} else {
		tc.log.Error("Empty Tickers", zap.String("Execution Level", "Controller"))
		c.JSON(http.StatusBadRequest, nil)
	}
}
