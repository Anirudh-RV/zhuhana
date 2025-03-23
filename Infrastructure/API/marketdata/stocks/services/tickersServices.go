package stocks

import (
	"marketdata/logger"
	tickerModels "marketdata/stocks/models"
	tickersRepository "marketdata/stocks/repositories"

	"go.uber.org/zap"
)

type TickersService struct {
	tickersRepository *tickersRepository.TickersRepository
	log               *logger.Logger
}

func NewTickersService(tickersRepository *tickersRepository.TickersRepository, log *logger.Logger) *TickersService {
	return &TickersService{
		tickersRepository: tickersRepository,
		log:               log,
	}
}

func (ts *TickersService) GetAllTickersV1(limit int) (*tickerModels.AllTickersAPIResponse, error) {
	tickers, err := ts.tickersRepository.GetAllTickersV1(limit)
	if err != nil && tickers == nil {
		ts.log.Error("Error getting Tickers", zap.String("Execution Level", "Service"), zap.String("Error", err.Error()))
		return nil, err
	}
	return tickers, err
}
