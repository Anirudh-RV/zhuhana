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

func (bts *BacktestService) GetStockDataWithRange(symbol, market string, from, to time.Time, limit, offset int) ([]models.OHLC, int, error) {
	ohlc_records, total, err := bts.backtestRepository.GetOHLCDataWithDateRange(symbol, market, from, to, limit, offset)
	if err != nil {
		go bts.logger.Error("error retrieving ohlc data with range", zap.String("execution level", "GetStockDataWithRange"), zap.String("Error", err.Error()))
		return nil, 0, err
	}
	return ohlc_records, int(total), nil
}

func (bts *BacktestService) GetStockDataAtTime(symbol, market string, current time.Time) (*models.OHLC, error) {
	ohlc_record, err := bts.backtestRepository.GetOHLCDataByTimestamp(symbol, market, current)
	if err != nil {
		go bts.logger.Error("error retrieving ohlc data with range", zap.String("execution level", "GetStockDataAtTime"), zap.String("Error", err.Error()))
		return nil, err
	}
	return ohlc_record, nil
}
