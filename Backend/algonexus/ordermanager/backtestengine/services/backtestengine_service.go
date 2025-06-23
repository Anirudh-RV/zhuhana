package services

import (
	"algonexus/logger"
	"algonexus/ordermanager/backtestengine/models"
	"algonexus/ordermanager/backtestengine/repositories"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"go.uber.org/zap"
)

type BacktestService struct {
	logger             *logger.Logger
	clickHouse         *clickhouse.Conn
	backtestRepository *repositories.BacktestRepository
}

func NewBacktestService(logger *logger.Logger, clickHouse *clickhouse.Conn, backtestRepository *repositories.BacktestRepository) *BacktestService {
	return &BacktestService{
		logger:             logger,
		clickHouse:         clickHouse,
		backtestRepository: backtestRepository,
	}
}

func (bts *BacktestService) GetStockData(domain string, symbol string, date string) error {
	return nil
}

func (bts *BacktestService) GetStockDataWithRange(symbol, market string, from, to time.Time) ([]models.OHLC, error) {
	ohlc_records, err := bts.backtestRepository.GetOHLCDataWithDateRange(symbol, market, from, to)
	if err != nil {
		go bts.logger.Error("error retreiving ohlc data with range", zap.String("execution level", "GetStockDataWithRange"), zap.String("Error", err.Error()))
		return nil, err
	}

	return ohlc_records, nil
}
