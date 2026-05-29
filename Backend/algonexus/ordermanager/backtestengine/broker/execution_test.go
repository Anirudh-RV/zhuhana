package broker

import (
	"algonexus/constants"
	"algonexus/logger"
	"algonexus/ordermanager/backtestengine/broker/services"
	"algonexus/ordermanager/models"
	"context"
	"errors"
	"testing"
	"time"
)

func testOrder(id string) models.OrderRequest {
	return models.OrderRequest{
		OrderID: id,
		Order: models.Order{
			Symbol:   "AAPL",
			Mode:     models.ModeIntraday,
			Side:     models.SideBuy,
			Type:     models.TypeMarket,
			Domain:   models.DomainBacktest,
			Quantity: 10,
			Price:    150.25,
		},
	}
}

// Happy path: a submitted order is executed by the worker pool and a BROKER_FILLED
// event comes back on Fills().
func TestInProcessBroker_HappyPathFill(t *testing.T) {
	b := NewInProcessBroker(logger.NewLogger(), services.NewMockSimulator(logger.NewLogger()), FailFast)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	b.Start(ctx)

	if err := b.Submit(ctx, testOrder("happy-1")); err != nil {
		t.Fatalf("Submit returned error: %v", err)
	}

	select {
	case evt := <-b.Fills():
		if evt.OrderID != "happy-1" {
			t.Errorf("OrderID = %q, want happy-1", evt.OrderID)
		}
		if evt.Status != models.StatusBrokerFilled {
			t.Errorf("Status = %q, want %q", evt.Status, models.StatusBrokerFilled)
		}
		if evt.Fill.Quantity != 10 {
			t.Errorf("Fill.Quantity = %v, want 10", evt.Fill.Quantity)
		}
		if evt.Fill.Price != 150.25 {
			t.Errorf("Fill.Price = %v, want 150.25", evt.Fill.Price)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for fill")
	}
}

// FailFast: when submitCh is full (no workers draining) Submit rejects immediately.
func TestInProcessBroker_FailFastWhenFull(t *testing.T) {
	b := NewInProcessBroker(logger.NewLogger(), services.NewMockSimulator(logger.NewLogger()), FailFast)
	// Intentionally do NOT Start() — nothing drains submitCh.
	ctx := context.Background()
	for i := 0; i < constants.SubmitBuffer; i++ {
		if err := b.Submit(ctx, testOrder("fill")); err != nil {
			t.Fatalf("Submit %d filled buffer unexpectedly early: %v", i, err)
		}
	}
	if err := b.Submit(ctx, testOrder("overflow")); !errors.Is(err, ErrBrokerBusy) {
		t.Fatalf("Submit on full queue = %v, want ErrBrokerBusy", err)
	}
}

// Backpressure: when submitCh is full Submit blocks rather than rejecting; here it
// blocks until the context deadline (proving it did not fail fast).
func TestInProcessBroker_BackpressureBlocks(t *testing.T) {
	b := NewInProcessBroker(logger.NewLogger(), services.NewMockSimulator(logger.NewLogger()), Backpressure)
	ctx := context.Background()
	for i := 0; i < constants.SubmitBuffer; i++ {
		if err := b.Submit(ctx, testOrder("fill")); err != nil {
			t.Fatalf("Submit %d filled buffer unexpectedly early: %v", i, err)
		}
	}
	cctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()
	start := time.Now()
	err := b.Submit(cctx, testOrder("blocked"))
	if err == nil {
		t.Fatal("Backpressure Submit on full queue returned nil, want it to block then return a ctx error")
	}
	if elapsed := time.Since(start); elapsed < 90*time.Millisecond {
		t.Fatalf("Submit returned after %v; expected to block ~100ms (did not actually block)", elapsed)
	}
}
