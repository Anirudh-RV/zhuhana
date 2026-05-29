package constants

// In-process order pipeline tunables.
//
// Only ONE Redis stream remains: OrderStream — the "submit" ingress (user script /
// HTTP producer -> OrderHub anchor). Everything past the anchor is in-process: the
// OrderHub<->Broker handoff is Go channels (see backtestengine/broker.InProcessBroker),
// not Redis. The old SubmitStream/FillStream Redis hops were removed; their decoupling
// is now provided by submitCh/fillCh inside the in-process broker.
const (
	// OrderStream is the order-system ingress stream consumed by the anchor.
	OrderStream         = "orderstream:strategy-1"
	OrderStreamGroup    = "group:order:strategy-1"
	OrderStreamConsumer = "order-strategy-1"

	// StreamMaxLen bounds the ingress stream length via approximate XADD trimming,
	// preventing the unbounded Redis growth the load test surfaced.
	StreamMaxLen = 100000

	// AnchorConcurrency caps in-flight orders the anchor (submit-worker pool) handles
	// concurrently. Orders are independent (per-order FSM) so no FIFO ordering is
	// required; a bounded worker pool replaces the old per-batch WaitGroup barrier.
	AnchorConcurrency = 16

	// BrokerConcurrency caps in-flight broker executions — the in-process broker's
	// worker pool draining submitCh.
	BrokerConcurrency = 16

	// FillsConcurrency caps the shared OrderHub Listener pool draining the broker's
	// fillCh (finalizes the FSM terminal tail + removes the handle from the registry).
	FillsConcurrency = 16

	// SubmitBuffer / FillBuffer size the in-process broker channels — the in-memory
	// analog of the removed SubmitStream/FillStream: bounded buffers that decouple
	// (potentially slow) execution from the order state machine. When submitCh is full
	// the broker applies its OverflowPolicy (fail-fast by default).
	SubmitBuffer = 1024
	FillBuffer   = 1024
)
