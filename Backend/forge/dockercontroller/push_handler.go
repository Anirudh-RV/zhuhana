package dockercontroller

import (
	"fmt"
	"os/exec"

	"go.uber.org/zap"
)

func (ds *DockerService) ImagePush(dockerImageName string) error {
	// Compose full image name (including username if needed)
	fullImageName := fmt.Sprintf("%s/%s", DOCKER_USERNAME, dockerImageName)
	ds.logger.Info("Pushing Docker image...",
		zap.String("image", fullImageName),
		zap.String("execution level", "ImagePush"))

	pushCmd := exec.Command("docker", "push", fullImageName)

	pushOut, err := pushCmd.CombinedOutput()
	if err != nil {
		ds.logger.Error("Docker push failed",
			zap.String("output", string(pushOut)),
			zap.Error(err))
		return err
	}
	ds.logger.Info("Docker image pushed successfully",
		zap.String("execution level", "ImagePush"),
		zap.String("output", string(pushOut)))

	return nil
}
