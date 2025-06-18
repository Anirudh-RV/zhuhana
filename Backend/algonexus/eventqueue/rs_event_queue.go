package eventqueue

import (
	"algonexus/logger"
	"context"
	"fmt"
	"go.uber.org/zap"
	"os"

	"github.com/redis/go-redis/v9"
)

type RedisStreamEventQueue struct {
	ctx    context.Context
	rdb    *redis.Client
	logger *logger.Logger
}

func NewRedisStreamEventQueue(ctx context.Context, logger *logger.Logger) *RedisStreamEventQueue {
	addr := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	password := os.Getenv("REDIS_PASSWORD")

	go logger.Info("Logging in to Redis ...",
		zap.String("registry", addr),
		zap.String("execution level", "EventQueue"))

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0, // Order Event Queue
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		logger.Fatal("Fail to connect to Redis", zap.Error(err),
			zap.String("registry", addr),
			zap.String("execution level", "EventQueue"))
		panic(err)
	}

	return &RedisStreamEventQueue{
		ctx:    ctx,
		rdb:    rdb,
		logger: logger,
	}
}
