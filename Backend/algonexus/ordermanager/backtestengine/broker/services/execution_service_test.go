package services

import (
	"algonexus/logger"
	"algonexus/ordermanager/models"
	"context"
	"testing"
)

func newOrderRequest(price, qty float64) models.OrderRequest {
	return models.OrderRequest{
		OrderID: "test-order-1",
		Order: models.Order{
			Symbol:   "SPY",
			Mode:     models.ModeIntraday,
			Side:     models.SideBuy,
			Type:     models.TypeMarket,
			Domain:   models.DomainBacktest,
			Quantity: qty,
			Price:    price,
		},
	}
}

func TestMockSimulator_FullFillAtOrderPrice(t *testing.T) {
	sim := NewMockSimulator(logger.NewLogger())

	if sim.Source() != "mock" || sim.Mode() != "simulator" {
		t.Fatalf("address mismatch: got %s/%s, want mock/simulator", sim.Source(), sim.Mode())
	}

	fill, err := sim.Execute(context.Background(), newOrderRequest(150.25, 10))
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	if fill.OrderID != "test-order-1" {
		t.Errorf("OrderID = %q, want test-order-1", fill.OrderID)
	}
	if fill.Price != 150.25 {
		t.Errorf("Price = %v, want 150.25 (echo order price)", fill.Price)
	}
	if fill.Quantity != 10 {
		t.Errorf("Quantity = %v, want 10 (full fill)", fill.Quantity)
	}
	if fill.Side != string(models.SideBuy) {
		t.Errorf("Side = %q, want BUY", fill.Side)
	}
	if fill.Status != string(models.StatusBrokerFilled) {
		t.Errorf("Status = %q, want %q", fill.Status, models.StatusBrokerFilled)
	}
	if fill.FillID == "" {
		t.Error("FillID is empty, want a uuid")
	}
}

func TestMockSimulator_FallbackPriceWhenUnset(t *testing.T) {
	sim := NewMockSimulator(logger.NewLogger())

	fill, err := sim.Execute(context.Background(), newOrderRequest(0, 5))
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
	if fill.Price != 100.0 {
		t.Errorf("Price = %v, want 100.0 fallback", fill.Price)
	}
}
