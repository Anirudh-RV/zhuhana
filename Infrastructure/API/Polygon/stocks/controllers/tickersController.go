package stocks

import (
	"fmt"
	"net/http"
	"polygon/logger"
	tickerModels "polygon/stocks/models"
	tickersService "polygon/stocks/services"
	"strconv"

	"github.com/gin-gonic/gin"
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

func (tc *TickersController) GetAllTickersV1(c *gin.Context) {
	var tickers *tickerModels.AllTickersAPIResponse
	var err error
	param := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(param)
	tickers, err = tc.tickersService.GetAllTickersV1(limit)

	if err != nil {
		fmt.Printf("ERROR - %s", err)
		c.JSON(http.StatusBadRequest, err)
	}
	if tickers != nil {
		c.JSON(http.StatusOK, tickers)
	} else {
		c.JSON(http.StatusBadRequest, nil)
	}
	return
}
