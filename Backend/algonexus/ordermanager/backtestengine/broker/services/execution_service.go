package services

import (
	"algonexus/logger"
	"algonexus/ordermanager/models"
	"context"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Executor produces a fill for an order. Implementations are addressed by
// (source, mode) — the data source they price against and their execution mode,
// e.g. ("mock", "simulator"). New brokers (backtest/OHLC, paper, live) plug in
// by implementing this interface and registering under their own (source, mode).
type Executor interface {
	Source() string
	Mode() string
	Execute(ctx context.Context, req models.OrderRequest) (models.OrderFill, error)
}

// MockSimulator is the toy executor: source "mock", mode "simulator".
// It ignores order type and market data — every order is treated as a MARKET
// order and instantly, fully filled at the order's own price (fallback 100 when
// unset). It exists to exercise the full order↔broker loop and FSM, not to model
// a real venue. A future backtest executor would price against OHLC instead.
type MockSimulator struct {
	logger *logger.Logger
}

func NewMockSimulator(logger *logger.Logger) *MockSimulator {
	return &MockSimulator{logger: logger}
}

func (m *MockSimulator) Source() string { return "mock" }
func (m *MockSimulator) Mode() string   { return "simulator" }

func (m *MockSimulator) Execute(ctx context.Context, req models.OrderRequest) (models.OrderFill, error) {
	price := req.Order.Price
	if price <= 0 {
		price = 100.0
	}

	fill := models.OrderFill{
		FillID:   uuid.New().String(),
		OrderID:  req.OrderID,
		Price:    price,
		Quantity: req.Order.Quantity, // MARKET: full fill
		Side:     string(req.Order.Side),
		Source:   m.Source() + "-" + m.Mode(),
		Time:     time.Now(),
		Status:   string(models.StatusBrokerFilled),
	}

	m.logger.Info("mock broker filled order",
		zap.String("orderID", req.OrderID),
		zap.String("symbol", req.Order.Symbol),
		zap.Float64("price", price),
		zap.Float64("qty", fill.Quantity),
		zap.String("execution level", "BrokerExecution mock/simulator"),
	)

	return fill, nil
}
