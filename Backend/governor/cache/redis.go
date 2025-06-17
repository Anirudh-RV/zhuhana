package cache

import (
	"context"
	"fmt"
	"governor/logger"
	"log"
	"os"

	"github.com/bsm/redislock"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var RedisObj *redis.Client
var RedisLockObj *redislock.Client

func InitRedis(ctx context.Context, logger *logger.Logger) {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)

	RedisObj = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
	})
	go logger.Info("redis instance creation successful", zap.String("Execution Level", "InitRedis"))

	errSet := RedisObj.Set(ctx, "key", "I got the value", 0).Err()
	if errSet != nil {
		log.Fatalf("error conneting to redis: %v", errSet)
	}

	value, errGet := RedisObj.Get(ctx, "key").Result()
	if errGet != nil {
		panic(errGet)
	}

	go logger.Info("redis SET/GET method tested. value: "+value, zap.String("Execution Level", "InitRedis"))

	RedisLockObj = redislock.New(RedisObj)
	go logger.Info("redis lock object initialized", zap.String("Execution Level", "InitRedis"))
}
