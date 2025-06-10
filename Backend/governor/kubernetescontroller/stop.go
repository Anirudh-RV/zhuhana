package kubernetescontroller

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (ks *KubernetesService) Stop(userAlgorithmID uuid.UUID) {
	// TEST THIS OUT
	userAlgorithmRunID, err := ks.GetUserAlgorithmRunByUserAlgorithmID(userAlgorithmID.String())
	if err != nil {
		fmt.Printf("Error detected %s\n", err.Error())
	}
	jobName := fmt.Sprintf("%s/%s", userAlgorithmID, userAlgorithmRunID)

	// Optional: Clean up Job and pod
	deletePolicy := meta.DeletePropagationForeground
	err = ks.clientSet.BatchV1().Jobs(ks.namespace).Delete(context.TODO(), jobName, meta.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
	if err != nil {
		fmt.Printf("Error detected %s\n", err.Error())
	}
	fmt.Printf("Job %s deleted\n", jobName)
}
