package scheduler

import (
	"database/sql"
	"governor/logger"

	"github.com/bsm/redislock"
	"github.com/robfig/cron/v3"
)

var RedisLockObj *redislock.Client
var Logger *logger.Logger
var DB *sql.DB
var CronScheduler *cron.Cron

func Init(redisLockObj *redislock.Client, logger *logger.Logger, db *sql.DB) {
	RedisLockObj = redisLockObj
	Logger = logger
	DB = db
	CronScheduler = cron.New()
	CronScheduler.Start()
	LoadCronJob()
}
