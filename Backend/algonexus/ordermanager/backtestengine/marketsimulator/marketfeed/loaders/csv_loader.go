package loaders

import (
	"algonexus/ordermanager/backtestengine/models"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

func LoadTicksFromCSV() ([]models.MarketTick, error) {
	path := "ordermanager/backtestengine/marketsimulator/marketfeed/repositories/data/csv/polygon_sample.csv"

	cwd, _ := os.Getwd()
	fmt.Println("Looking for:", path, "in", cwd)

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	_, err = reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	var ticks []models.MarketTick
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, fmt.Errorf("failed to read record: %w", err)
		}

		volume, _ := strconv.Atoi(record[1])
		open, _ := strconv.ParseFloat(record[2], 64)
		closePrice, _ := strconv.ParseFloat(record[3], 64)
		high, _ := strconv.ParseFloat(record[4], 64)
		low, _ := strconv.ParseFloat(record[5], 64)

		startNs, _ := strconv.ParseInt(record[6], 10, 64)
		startTime := time.Unix(0, startNs)

		transactions, _ := strconv.ParseFloat(record[7], 64)

		tick := models.MarketTick{
			Ticker:       record[0],
			Volume:       volume,
			Open:         open,
			Close:        closePrice,
			High:         high,
			Low:          low,
			Start:        startTime,
			Transactions: transactions,
		}

		ticks = append(ticks, tick)
	}

	return ticks, nil
}
