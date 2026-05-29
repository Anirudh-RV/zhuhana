package broker

import (
	"algonexus/constants"
	"algonexus/logger"
	brokermodels "algonexus/ordermanager/backtestengine/broker/models"
	"algonexus/ordermanager/backtestengine/broker/services"
	ordermodels "algonexus/ordermanager/models"
	"context"
	"errors"
	"strings"

	"go.uber.org/zap"
)

// ErrBrokerBusy is returned by Submit under the FailFast policy when submitCh is full.
var ErrBrokerBusy = errors.New("broker submit queue full")

// OverflowPolicy decides what Submit does when the submit queue is full.
type OverflowPolicy int

const (
	// FailFast rejects immediately when submitCh is full (default). Saturation is
	// visible (the order goes ERROR) and pinpoints the bottleneck fast. See
	// loadtest/BACKPRESSURE.md.
	FailFast OverflowPolicy = iota
	// Backpressure blocks until the queue drains. NOTE: incomplete without an ingress
	// admission gate — on its own it merely relocates the pile-up to orderstream
	// (which then trims via MAXLEN). Kept as a switchable placeholder.
	Backpressure
)

// ParseOverflowPolicy reads the policy from env (BROKER_OVERFLOW_POLICY). Default FailFast.
func ParseOverflowPolicy(s string) OverflowPolicy {
	if strings.EqualFold(strings.TrimSpace(s), "backpressure") {
		return Backpressure
	}
	return FailFast
}

// InProcessBroker is the backtest-engine broker as an in-process module. OrderHub talks
// to it ONLY through the Broker port (Submit + Fills) — there is no Redis between them.
// submitCh/fillCh are the in-memory analog of the removed SubmitStream/FillStream:
// bounded buffers that decouple (potentially slow) execution from the order FSM.
type InProcessBroker struct {
	logger   *logger.Logger
	executor services.Executor
	policy   OverflowPolicy
	submitCh chan ordermodels.OrderRequest
	fillCh   chan brokermodels.BrokerFillEvent
	workers  int
}

func NewInProcessBroker(log *logger.Logger, executor services.Executor, policy OverflowPolicy) *InProcessBroker {
	return &InProcessBroker{
		logger:   log,
		executor: executor,
		policy:   policy,
		submitCh: make(chan ordermodels.OrderRequest, constants.SubmitBuffer),
		fillCh:   make(chan brokermodels.BrokerFillEvent, constants.FillBuffer),
		workers:  constants.BrokerConcurrency,
	}
}

// Start launches the bounded worker pool that drains submitCh.
func (b *InProcessBroker) Start(ctx context.Context) {
	b.logger.Info("InProcessBroker starting",
		zap.Int("workers", b.workers),
		zap.String("policy", b.policyName()),
		zap.String("execution level", "InProcessBroker"))
	for i := 0; i < b.workers; i++ {
		go b.worker(ctx)
	}
}

func (b *InProcessBroker) worker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case req := <-b.submitCh:
			evt := brokermodels.BrokerFillEvent{OrderID: req.OrderID}
			fill, err := b.executor.Execute(ctx, req)
			if err != nil {
				evt.Status = ordermodels.StatusBrokerRejected
				b.logger.Error("broker execute failed",
					zap.String("orderID", req.OrderID), zap.Error(err),
					zap.String("execution level", "InProcessBroker"))
			} else {
				evt.Fill = fill
				evt.Status = ordermodels.StatusBrokerFilled
			}
			select {
			case b.fillCh <- evt:
			case <-ctx.Done():
				return
			}
		}
	}
}

// Submit hands an order to the broker (command). Returns nil on accept. Under FailFast
// it returns ErrBrokerBusy when submitCh is full; under Backpressure it blocks until a
// slot frees (or ctx is cancelled).
func (b *InProcessBroker) Submit(ctx context.Context, req ordermodels.OrderRequest) error {
	if b.policy == Backpressure {
		select {
		case b.submitCh <- req:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	// FailFast: never block.
	select {
	case b.submitCh <- req:
		return nil
	default:
		return ErrBrokerBusy
	}
}

// Fills is the broker's async execution-result stream (event), consumed by the OrderHub
// Listener pool.
func (b *InProcessBroker) Fills() <-chan brokermodels.BrokerFillEvent {
	return b.fillCh
}

func (b *InProcessBroker) policyName() string {
	if b.policy == Backpressure {
		return "backpressure"
	}
	return "failfast"
}
