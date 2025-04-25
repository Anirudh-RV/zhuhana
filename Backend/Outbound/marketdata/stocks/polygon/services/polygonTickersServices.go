package polygon

import (
	"outbound/logger"
	tickerModels "outbound/marketdata/stocks/models"
	tickersRepository "outbound/marketdata/stocks/polygon/repositories"
)

type PolygonTickersService struct {
	polygonTickersRepository *tickersRepository.PolygonTickersRepository
	log                      *logger.Logger
}

func NewPolygonTickersService(tickersRepository *tickersRepository.PolygonTickersRepository, log *logger.Logger) *PolygonTickersService {
	return &PolygonTickersService{
		polygonTickersRepository: tickersRepository,
		log:                      log,
	}
}

func (pts *PolygonTickersService) GetDailyTickerOHLCV_V1(limit int) (*tickerModels.DailyTickerOHLCVResponse, error) {

}
