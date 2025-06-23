package controllers

import (
	"algonexus/constants"
	"algonexus/ordermanager/backtestengine/models"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GetOHLCDataWithNext godoc
//
// @Summary      Get OHLC data for a symbol at a specific timestamp
// @Description  Retrieves OHLC (Open, High, Low, Close) data for a given symbol and market at the specified timestamp (`current_time`). Also returns a `next_url` link with `next_step` seconds added to `current_time` if within `end_time`.
// @Tags         Backtest
// @Accept       json
// @Produce      json
// @Param        current_time  query     string  true  "Current timestamp in RFC3339 format"     example(2020-06-10T00:00:00Z)
// @Param        end_time      query     string  true  "End timestamp in RFC3339 format"         example(2020-07-15T00:00:00Z)
// @Param        symbol        query     string  true  "Symbol to retrieve OHLC data for"        example(SPY)
// @Param        market        query     string  true  "Market for the symbol"                   example(NYSEARCA)
// @Param        next_step     query     int     true  "Step in seconds to calculate next_url"   example(86400)
// @Success      200  {object}  models.OHLCNextResponse
// @Failure      400  {object}  models.OHLCNextResponse  "Invalid request format or query params"
// @Failure      500  {object}  models.OHLCNextResponse  "Internal server error while retrieving data"
// @Security     ApiKeyAuth
// @Router       /v1/backtest/ohlc/next/ [get]
func (btc *BacktestController) GetOHLCDataWithNext(c *gin.Context) {
	currentStr := c.Query("current_time")
	endStr := c.Query("end_time")
	symbol := c.Query("symbol")
	market := c.Query("market")
	nextStepStr := c.Query("next_step")

	current, err := time.Parse(time.RFC3339, currentStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.OHLCNextResponse{
			Status:            0,
			StatusDescription: "Invalid start_time format, must be RFC3339 with timezone",
		})
		return
	}
	end, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.OHLCNextResponse{
			Status:            0,
			StatusDescription: "Invalid end_time format, must be RFC3339 with timezone",
		})
		return
	}

	nextStep, err := strconv.Atoi(nextStepStr)
	if err != nil || nextStep <= 0 {
		c.JSON(http.StatusBadRequest, models.OHLCNextResponse{
			Status:            0,
			StatusDescription: "Invalid next_step, must be a positive integer (seconds)",
		})
		return
	}

	ohlcData, err := btc.backtestService.GetStockDataAtTime(symbol, market, current.UTC())
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.OHLCNextResponse{
			Status:            0,
			StatusDescription: "Server error while retrieving OHLC data",
		})
		return
	}

	var nextURL *string
	nextTime := current.Add(time.Duration(nextStep) * time.Second)
	if nextTime.Before(end) {
		url := fmt.Sprintf(
			constants.ALGONEXUS_URL+"/v1/backtest/ohlc/next/?current_time=%s&end_time=%s&symbol=%s&market=%s&next_step=%d",
			nextTime.Format(time.RFC3339), end.Format(time.RFC3339), symbol, market, nextStep,
		)
		nextURL = &url
	}

	c.JSON(http.StatusOK, models.OHLCNextResponse{
		Status:            1,
		StatusDescription: "Successfully retrieved OHLC data",
		OHLCData:          ohlcData,
		NextURL:           nextURL,
	})
}
