package kubernetescontroller

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

func (ks *KubernetesService) Start(userAlgorithmID uuid.UUID) {
	// TEST THIS OUT
	userAlgorithm, err := ks.GetUserAlgorithm(userAlgorithmID)
	if err != nil {
		go ks.logger.Error("could not get useralgorithm from db", zap.String("Execution Level", "KubernetesStart"))
	}
	dockerImageName := fmt.Sprintf("user-algorithm-%s-%s", userAlgorithm.UserID, userAlgorithmID)
	pullImageName := fmt.Sprintf("%s/%s", DOCKER_REPOSITORY, dockerImageName)
	go ks.logger.Info(fmt.Sprintf("trying to start: %s \n full name: %s", dockerImageName, pullImageName), zap.String("Execution Level", "KubernetesStart"))

	userAlgorithmRunUUID, err := ks.AddUserAlgorithmRun(userAlgorithmID, *userAlgorithm.StartCronSchedule, *userAlgorithm.EndCronSchedule, int(userAlgorithm.OrderDomain))
	if err != nil {
		fmt.Printf("Error detected %s\n", err.Error())
	}
	userAlgorithmRunID := fmt.Sprint(userAlgorithmRunUUID)
	containerName := fmt.Sprintf("%s-%s", userAlgorithmID, userAlgorithmRunID)
	jobName := fmt.Sprintf("%s-%s", userAlgorithmID, userAlgorithmRunID)

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
							Args:  []string{}, // Optional: pass args here
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
		if err != nil || len(pods.Items) == 0 {
			return false, nil
		}
		pod := pods.Items[0]
		if pod.Status.Phase == corev1.PodSucceeded || pod.Status.Phase == corev1.PodFailed {
			podName = pod.Name
			return true, nil
		}
		return false, nil
	})
	if err != nil {
		panic(fmt.Sprintf("Job pod not ready in time: %v", err))
	}

	// Stream Logs
	req := ks.clientSet.CoreV1().Pods(ks.namespace).GetLogs(podName, &corev1.PodLogOptions{})
	stream, err := req.Stream(context.TODO())
	if err != nil {
		panic(err)
	}
	defer stream.Close()

	fmt.Println("Logs from job:")
	io.Copy(os.Stdout, stream) // Replace io.Discard with os.Stdout if you want to print logs
}
