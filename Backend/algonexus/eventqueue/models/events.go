package models

type EventType string

const (
	EventBrokerConfirmed EventType = "broker_confirmed"
	EventFill            EventType = "fill"
	EventComplete        EventType = "complete"
	EventError           EventType = "error"
	EventHeartbeat       EventType = "heartbeat" // optional
)

type RsEvent struct {
	OrderID   string    `json:"order_id"`
	Timestamp int64     `json:"timestamp"` // optional: Unix time
	Type      EventType `json:"type"`
	Payload   any       `json:"payload,omitempty"` // extra data if needed
}
