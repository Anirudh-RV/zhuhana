package kubernetescontroller

import (
	"context"
	"fmt"
	"governor/constants"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

func (ks *KubernetesService) Start(userAlgorithmID uuid.UUID, market, symbol string, startTime, endTime *time.Time, portfolioSize, frequency int) error {
	userAlgorithm, err := ks.GetUserAlgorithm(userAlgorithmID)
	if err != nil {
		go ks.logger.Error("could not get useralgorithm from db", zap.String("Execution Level", "KubernetesStart"))
		return err
	}
	userAlgorithmToken, err := ks.GetUserAlgorithmToken(userAlgorithmID.String())
	if err != nil {
		go ks.logger.Error("could not get userAlgorithmToken", zap.String("Execution Level", "KubernetesStart"), zap.String("Error", err.Error()))
		return err
	}
	if userAlgorithmToken == "" {
		go ks.logger.Error("user algorithm token is empty", zap.String("Execution Level", "KubernetesStart"))
		return fmt.Errorf("user algorithm token is empty")
	}

	dockerImageName := fmt.Sprintf("user-algorithm-%s-%s", userAlgorithm.UserID, userAlgorithmID)
	pullImageName := fmt.Sprintf("%s/%s", DOCKER_REPOSITORY, dockerImageName)
	go ks.logger.Info(fmt.Sprintf("trying to start: %s \n full name: %s", dockerImageName, pullImageName), zap.String("Execution Level", "KubernetesStart"))

	startCronSchedule := ""
	if userAlgorithm.StartCronSchedule != nil {
		startCronSchedule = *userAlgorithm.StartCronSchedule
	}

	endCronSchedule := ""
	if userAlgorithm.EndCronSchedule != nil {
		endCronSchedule = *userAlgorithm.EndCronSchedule
	}
	userAlgorithmRunUUID, err := ks.AddUserAlgorithmRun(
		userAlgorithmID,
		startCronSchedule,
		endCronSchedule,
		int(userAlgorithm.OrderDomain),
		market,
		symbol,
		startTime,
		endTime,
		portfolioSize)
	if err != nil {
		fmt.Printf("Error detected %s\n", err.Error())
		return err
	}
	userAlgorithmRunID := fmt.Sprint(userAlgorithmRunUUID)
	containerName := userAlgorithmRunID
	jobName := userAlgorithmRunID

	// Define Job
	job := &batchv1.Job{
		ObjectMeta: meta.ObjectMeta{
			Name: jobName,
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyNever,
					Containers: []corev1.Container{
						{
							Name:  containerName,
							Image: pullImageName,
							Args:  []string{},
							Env: []corev1.EnvVar{
								{
									Name:  "USER_ALGORITHM_TOKEN",
									Value: userAlgorithmToken,
								},
								{
									Name:  "ORDER_DOMAIN",
									Value: userAlgorithm.OrderDomain.String(),
								},
								{
									Name:  "MARKET",
									Value: market,
								},
								{
									Name:  "SYMBOL",
									Value: symbol,
								},
								{
									Name:  "START_TIME",
									Value: startTime.String(),
								},
								{
									Name:  "END_TIME",
									Value: endTime.String(),
								},
								{
									Name:  "PORTFOLIO_SIZE",
									Value: fmt.Sprint(portfolioSize),
								},
								{
									Name:  "FREQUENCY",
									Value: fmt.Sprint(frequency),
								},
								{
									Name:  "API_ENDPOINT",
									Value: constants.USER_ALGORITHM_API_ENDPOINT,
								},
							},
						},
					},
				},
			},
		},
	}

	// Create Job
	_, err = ks.clientSet.BatchV1().Jobs(ks.namespace).Create(context.TODO(), job, meta.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Job %s created\n", jobName)

	// Wait for Pod to be ready and completed
	var podName string
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	err = wait.PollUntilContextTimeout(ctx, 2*time.Second, 2*time.Minute, true, func(ctx context.Context) (bool, error) {
		pods, err := ks.clientSet.CoreV1().Pods(ks.namespace).List(ctx, meta.ListOptions{
			LabelSelector: "job-name=" + jobName,
		})
		if err != nil {
			fmt.Printf("[PollLoop] Error listing pods: %s\n", err)
			return false, nil
		}
		if len(pods.Items) == 0 {
			fmt.Println("[PollLoop] No pods found for job yet")
			return false, nil
		}

		for _, pod := range pods.Items {
			fmt.Printf("[PollLoop] Found pod: %s, Phase: %s\n", pod.Name, pod.Status.Phase)
			for _, cs := range pod.Status.ContainerStatuses {
				fmt.Printf("[PollLoop] Container %s - Ready: %v, State: %+v\n", cs.Name, cs.Ready, cs.State)
			}
		}

		pod := pods.Items[0] // Assuming first one
		switch pod.Status.Phase {
		case corev1.PodRunning:
			podName = pod.Name
			fmt.Printf("[PollLoop] Pod %s is now Running\n", podName)
			return true, nil
		case corev1.PodFailed:
			return false, fmt.Errorf("pod %s failed", pod.Name)
		default:
			fmt.Printf("[PollLoop] Pod %s in Phase: %s\n", pod.Name, pod.Status.Phase)
			return false, nil
		}
	})

	if err != nil {
		fmt.Printf("unable to stream logs: %s", err.Error())
	}

	// Stream Logs
	go func() {
		if err := ks.StreamPodLogs(podName); err != nil {
			fmt.Printf("log stream error: %s\n", err)
		}
	}()

	return nil
}
