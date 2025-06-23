package controllers

import (
	"algonexus/logger"
	"algonexus/ordermanager/backtestengine/models"
	"algonexus/ordermanager/backtestengine/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type BacktestController struct {
	backtestService *services.BacktestService
	log             *logger.Logger
}

func NewBacktestController(backtestService *services.BacktestService, log *logger.Logger) *BacktestController {
	return &BacktestController{
		backtestService: backtestService,
		log:             log,
	}
}

// GetOHLCDataWithRange godoc
//
// @Summary      Get OHLC data for a symbol within a date time range
// @Description  Retrieves OHLC (Open, High, Low, Close) data for a given symbol and market between the specified start and end times. Supports pagination.
// @Tags         Backtest
// @Accept       json
// @Produce      json
// @Param        start_time   query     string  true  "Start time in RFC3339 format"  example(2020-06-10T00:00:00Z)
// @Param        end_time     query     string  true  "End time in RFC3339 format"    example(2020-06-15T00:00:00Z)
// @Param        symbol       query     string  true  "Symbol to retrieve OHLC data for"  example(SPY)
// @Param        market       query     string  true  "Market for the symbol"             example(NYSEARCA)
// @Param        page_limit   query     int     false "Number of results per page"        default(50) minimum(1)
// @Param        page_no      query     int     false "Page number to retrieve"           default(1) minimum(1)
// @Success      200  {object}  models.OHLCRangeResponse
// @Failure      400  {object}  models.OHLCRangeResponse  "Invalid time format or missing required query parameters"
// @Failure      500  {object}  models.OHLCRangeResponse  "Internal server error"
// @Security     ApiKeyAuth
// @Router       /v1/backtest/ohlc/range/ [get]
func (btc *BacktestController) GetOHLCDataWithRange(c *gin.Context) {
	startStr := c.Query("start_time")
	endStr := c.Query("end_time")
	symbol := c.Query("symbol")
	market := c.Query("market")

	pageLimitStr := c.DefaultQuery("page_limit", "50")
	pageNoStr := c.DefaultQuery("page_no", "1")

	// Parse time range
	start, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.OHLCRangeResponse{
			Status:            0,
			StatusDescription: "Invalid start_time format, must be RFC3339 with timezone",
		})
		return
	}
	end, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.OHLCRangeResponse{
			Status:            0,
			StatusDescription: "Invalid end_time format, must be RFC3339 with timezone",
		})
		return
	}

	pageLimit, err := strconv.Atoi(pageLimitStr)
	if err != nil || pageLimit <= 0 {
		pageLimit = 50
	}
	pageNo, err := strconv.Atoi(pageNoStr)
	if err != nil || pageNo <= 0 {
		pageNo = 1
	}
	offset := (pageNo - 1) * pageLimit

	ohlcData, total, err := btc.backtestService.GetStockDataWithRange(symbol, market, start.UTC(), end.UTC(), pageLimit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.OHLCRangeResponse{
			Status:            0,
			StatusDescription: "Server error while retrieving OHLC data",
		})
		return
	}

	var nextPage *int
	if offset+pageLimit < total {
		n := pageNo + 1
		nextPage = &n
	}

	c.JSON(http.StatusOK, models.OHLCRangeResponse{
		Status:            1,
		StatusDescription: "Successfully retrieved OHLC data",
		OHLCData:          ohlcData,
		PageLimit:         pageLimit,
		PageNo:            pageNo,
		NextPageNo:        nextPage,
	})
}
