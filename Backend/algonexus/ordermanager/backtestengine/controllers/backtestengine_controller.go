package controllers

import (
	"algonexus/logger"
	"algonexus/ordermanager/backtestengine/services"
)

type BacktestController struct {
	backtestService *services.BacktestService
	log             *logger.Logger
}

func NewBacktestController(backtestService *services.BacktestService, log *logger.Logger) *BacktestController {
	return &BacktestController{
		backtestService: backtestService,
		log:             log,
	}
}
