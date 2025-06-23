package db

import (
	"algonexus/logger"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	clickhouse "github.com/ClickHouse/clickhouse-go/v2"
)

var ClickHouse clickhouse.Conn

func InitClickHouse(logger *logger.Logger) {
	host := os.Getenv("CLICKHOUSE_HOST")
	port := os.Getenv("CLICKHOUSE_PORT")
	user := os.Getenv("CLICKHOUSE_USER")
	password := os.Getenv("CLICKHOUSE_PASSWORD")
	dbname := os.Getenv("CLICKHOUSE_DB")

	addr := fmt.Sprintf("%s:%s", host, port)

	var err error
	ClickHouse, err = clickhouse.Open(&clickhouse.Options{
		Addr: []string{addr},
		Auth: clickhouse.Auth{
			Database: dbname,
			Username: user,
			Password: password,
		},
		DialTimeout: 5 * time.Second,
		Compression: &clickhouse.Compression{Method: clickhouse.CompressionLZ4},
		ClientInfo:  clickhouse.ClientInfo{},
		Debug:       false,
		Protocol:    clickhouse.Native, // or clickhouse.HTTP
	})

	if err != nil {
		log.Fatalf("ClickHouse connection init failed: %v", err)
	}

	// Ping to check if connection is alive
	ctx := context.Background()
	if err := ClickHouse.Ping(ctx); err != nil {
		log.Fatalf("Cannot connect to ClickHouse: %v", err)
	}
}
