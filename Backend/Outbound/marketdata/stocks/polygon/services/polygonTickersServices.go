package polygon

import (
	"outbound/logger"
	tickerModels "outbound/marketdata/stocks/models"
	tickersRepository "outbound/marketdata/stocks/polygon/repositories"

	"go.uber.org/zap"
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

func (pts *PolygonTickersService) GetAllTickersV1(limit int) (*tickerModels.AllTickersAPIResponse, error) {
	tickers, err := pts.polygonTickersRepository.GetAllTickersV1(limit)
	if err != nil && tickers == nil {
		go pts.log.Error("error getting Tickers", zap.String("execution level", "Service"), zap.String("error", err.Error()))
		return nil, err
	}
	return tickers, err
}
