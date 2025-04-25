package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"orchestrator/logger"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

var DB *sql.DB

func InitDB(logger *logger.Logger) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	go logger.Info("DSN for connection: "+dsn, zap.String("Execution Level", "Root"))

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("Cannot connect to DB: %v", err)
	}

	log.Println("Successfully connected to the database")
}
