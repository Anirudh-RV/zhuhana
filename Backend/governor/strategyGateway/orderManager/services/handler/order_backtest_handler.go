package handler

import (
	"governor/strategyGateway/orderManager/models"
	"time"
)

func HandleBacktestOrder(req *models.OrderRequest) (*models.OrderResponse, error) {
	res := &models.OrderResponse{
		OrderID:       req.OrderID,
		BrokerOrderID: "Zhuhana-backtest",
		Status:        models.ResponseStatusFilled,
		Message:       "Trade Success",
		Time:          time.Now(),
	}

	return res, nil
}
