package scheduler

import (
	"database/sql"
	"governor/kafka"
	"governor/logger"

	"github.com/bsm/redislock"
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
)

type SchedulerService struct {
	redisLockObj  *redislock.Client
	redisObj      *redis.Client
	logger        *logger.Logger
	db            *sql.DB
	cronScheduler *cron.Cron
	kafkaService  *kafka.KafkaService
}

func NewSchedulerService(redisObj *redis.Client, redisLockObj *redislock.Client, logger *logger.Logger, db *sql.DB, kafkaService *kafka.KafkaService) *SchedulerService {
	cronScheduler := cron.New()
	cronScheduler.Start()
	return &SchedulerService{
		redisObj:      redisObj,
		redisLockObj:  redisLockObj,
		logger:        logger,
		db:            db,
		cronScheduler: cronScheduler,
		kafkaService:  kafkaService,
	}
}
