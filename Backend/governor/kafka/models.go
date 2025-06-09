package kafka

import "time"

type JobPayload struct {
	JobID   string      `json:"job_id"`
	Target  string      `json:"target"`  // who should consume it, e.g., "governor"
	Payload interface{} `json:"payload"` // job data
	Time    time.Time   `json:"time"`    // time of creation
}
