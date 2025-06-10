package scheduler

import (
	"database/sql"
	"governor/kafka"
	"governor/logger"

	"github.com/bsm/redislock"
	"github.com/robfig/cron/v3"
)

type SchedulerService struct {
	redisLockObj  *redislock.Client
	logger        *logger.Logger
	db            *sql.DB
	cronScheduler *cron.Cron
	kafkaService  *kafka.KafkaService
}

func NewSchedulerService(redisLockObj *redislock.Client, logger *logger.Logger, db *sql.DB, kafkaService *kafka.KafkaService) *SchedulerService {
	cronScheduler := cron.New()
	cronScheduler.Start()
	return &SchedulerService{
		redisLockObj:  redisLockObj,
		logger:        logger,
		db:            db,
		cronScheduler: cronScheduler,
		kafkaService:  kafkaService,
	}
}
