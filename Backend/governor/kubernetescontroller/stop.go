package kubernetescontroller

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (ks *KubernetesService) Stop(userAlgorithmID uuid.UUID) error {
	// TEST THIS OUT
	userAlgorithmRunIDs, err := ks.GetUserAlgorithmRunsByUserAlgorithmID(userAlgorithmID.String())
	if err != nil {
		fmt.Printf("Error detected %s\n", err.Error())
		return err
	}

	fmt.Println("Found active user_algorithm_run IDs:")
	for _, id := range userAlgorithmRunIDs {
		fmt.Println("-", id)
	}

	for _, userAlgorithmRunID := range userAlgorithmRunIDs {
		jobName := userAlgorithmRunID.String()

		// Optional: Clean up Job and pod
		deletePolicy := meta.DeletePropagationForeground
		err = ks.clientSet.BatchV1().Jobs(ks.namespace).Delete(context.TODO(), jobName, meta.DeleteOptions{
			PropagationPolicy: &deletePolicy,
		})
		if err != nil {
			fmt.Printf("Error detected %s\n", err.Error())
		}
		fmt.Printf("Job %s deleted\n", jobName)
		ks.DeactivateUserAlgorithmRunByID(userAlgorithmRunID)
	}

	return nil
}
