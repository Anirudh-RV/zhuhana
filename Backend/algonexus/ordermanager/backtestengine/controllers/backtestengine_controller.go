package controllers

import (
	"algonexus/logger"
	"algonexus/ordermanager/backtestengine/models"
	"algonexus/ordermanager/backtestengine/services"
	"net/http"
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

func (btc *BacktestController) GetOHLCDataWithRange(c *gin.Context) {
	startStr := c.Query("start_time")
	endStr := c.Query("end_time")
	symbol := c.Query("symbol")
	market := c.Query("market")

	// Parse with timezone awareness
	start, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.OHLCRangeResponse{
			Status:            0,
			StatusDescription: "Invalid start_time format, must be RFC3339 with timezone",
		})
		return
	}
	start = start.UTC()

	end, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.OHLCRangeResponse{
			Status:            0,
			StatusDescription: "Invalid end_time format, must be RFC3339 with timezone",
		})
		return
	}
	end = end.UTC()

	ohlcData, err := btc.backtestService.GetStockDataWithRange(symbol, market, start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.OHLCRangeResponse{
			Status:            0,
			StatusDescription: "Server error while retrieving OHLC data",
		})
		return
	}

	c.JSON(http.StatusOK, models.OHLCRangeResponse{
		Status:            1,
		StatusDescription: "Successfully retrieved OHLC data",
		OHLCData:          ohlcData,
	})
}
