package kubernetescontroller

import (
	"bufio"
	"context"
	"fmt"
	"strings"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
)

func (ks *KubernetesService) StreamPodLogs(podName string) error {
	req := ks.clientSet.CoreV1().Pods(ks.namespace).GetLogs(podName, &corev1.PodLogOptions{
		Follow: true,
	})

	stream, err := req.Stream(context.TODO())
	if err != nil {
		return fmt.Errorf("failed to open log stream: %w", err)
	}
	defer stream.Close()

	scanner := bufio.NewScanner(stream)
	for scanner.Scan() {
		line := scanner.Text()
		go ks.logger.Info(fmt.Sprintf("Log: %s", line), zap.String("Execution Level", "KubernetesStart"))

		// 🧠 Do custom processing here:
		if strings.Contains(line, "error") {
			// Log to your logger
			ks.logger.Error("Error found in logs", zap.String("log", line))
		}

		if strings.Contains(line, "Job Completed") {
			// Optionally trigger some follow-up action
			go ks.logger.Info("job completed", zap.String("Execution Level", "KubernetesStart"))
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading log stream: %w", err)
	}

	return nil
}
