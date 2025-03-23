package stocks

import (
	"polygon/logger"
	tickerModels "polygon/stocks/models"
	tickersRepository "polygon/stocks/repositories"
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
		return nil, err
	}
	return tickers, err
}
