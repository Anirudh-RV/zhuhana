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

func (ur *BacktestRepository) GetOHLCDataWithDateRange(symbol, market string, from, to time.Time, limit, offset int) ([]models.OHLC, uint64, error) {
	ctx := context.Background()

	// Get total count
	countQuery := `
		SELECT count(*)
		FROM OHLC
		WHERE Symbol = ? AND Market = ? AND Date_Time >= ? AND Date_Time <= ?
	`
	var total uint64
	if err := db.ClickHouse.QueryRow(ctx, countQuery, symbol, market, from, to).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Fetch paginated results
	dataQuery := `
		SELECT Symbol, Market, Date_Time, Open, High, Low, Close, Volume, Day, Weekday, Week, Month, Year
		FROM OHLC
		WHERE Symbol = ? AND Market = ? AND Date_Time >= ? AND Date_Time <= ?
		ORDER BY Date_Time
		LIMIT ? OFFSET ?
	`

	rows, err := db.ClickHouse.Query(ctx, dataQuery, symbol, market, from, to, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var result []models.OHLC
	for rows.Next() {
		var row models.OHLC
		if err := rows.ScanStruct(&row); err != nil {
			return nil, 0, err
		}
		result = append(result, row)
	}

	return result, total, nil
}

func (ur *BacktestRepository) GetOHLCDataByTimestamp(symbol, market string, timestamp time.Time) (*models.OHLC, error) {
	ctx := context.Background()

	query := `
		SELECT Symbol, Market, Date_Time, Open, High, Low, Close, Volume, Day, Weekday, Week, Month, Year
		FROM OHLC
		WHERE Symbol = ? AND Market = ? AND Date_Time = ?
		LIMIT 1
	`

	row := db.ClickHouse.QueryRow(ctx, query, symbol, market, timestamp)

	var result models.OHLC
	if err := row.ScanStruct(&result); err != nil {
		// You might want to use sql.ErrNoRows check here if needed
		return nil, err
	}

	return &result, nil
}
