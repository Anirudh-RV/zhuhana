package kafka

import "log"

func KafkaConsumer(job JobPayload) error {
	// TODO: Get the kubernetes to run the container
	log.Printf("Received job %s", job.JobID)
	return nil
}
