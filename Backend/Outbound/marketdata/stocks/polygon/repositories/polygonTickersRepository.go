package polygon

import (
	"outbound/logger"
	tickerModels "outbound/marketdata/stocks/models"
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
	GetDailyTickerOHLCV_V1(limit int) (*tickerModels.DailyTickerOHLCVResponse, error)
}

func (ptr *PolygonTickersRepository) GetDailyTickerOHLCV_V1(limit int) (*tickerModels.DailyTickerOHLCVResponse, error) {

}
