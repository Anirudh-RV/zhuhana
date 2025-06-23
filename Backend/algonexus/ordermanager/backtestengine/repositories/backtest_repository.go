package repositories

import (
	"algonexus/db"
	"algonexus/ordermanager/backtestengine/models"
	"context"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
)

type BacktestRepository struct {
	clickHouse *clickhouse.Conn
}

func NewBacktestRepository(clickHouse *clickhouse.Conn) *BacktestRepository {
	return &BacktestRepository{
		clickHouse: clickHouse,
	}
}

func (ur *BacktestRepository) GetOHLCDataWithDateRange(symbol, market string, from, to time.Time) ([]models.OHLC, error) {
	ctx := context.Background()
	query := `
		SELECT Symbol, Market, Date_Time, Open, High, Low, Close, Volume, Day, Weekday, Week, Month, Year
		FROM OHLC
		WHERE Symbol = ? AND Market = ? AND Date_Time >= ? AND Date_Time <= ?
		ORDER BY Date_Time
	`

	rows, err := db.ClickHouse.Query(ctx, query, symbol, market, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.OHLC
	for rows.Next() {
		var row models.OHLC
		if err := rows.ScanStruct(&row); err != nil {
			return nil, err
		}
		result = append(result, row)
	}

	return result, nil
}
