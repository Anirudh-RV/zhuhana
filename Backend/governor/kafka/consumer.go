package kafka

import "log"

func KafkaConsumer(job JobPayload) error {
	// TODO: Get the kubernetes to run the container
	log.Printf("kafka consumer received job %s", job)
	return nil
}
