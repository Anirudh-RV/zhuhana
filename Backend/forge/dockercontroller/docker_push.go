package dockercontroller

import (
	"fmt"
	"os/exec"

	"go.uber.org/zap"
)

func (ds *DockerService) ImagePush(dockerImageName string) error {
	// Compose full image name (including username if needed)
	fullImageName := fmt.Sprintf("%s/%s", DOCKER_USERNAME, dockerImageName)

	// Step 1: Login to Docker registry
	ds.logger.Info("Logging in to Docker registry...",
		zap.String("registry", DOCKER_SERVER_ADDRESS),
		zap.String("execution level", "ImagePush"))

	loginCmd := exec.Command("docker", "login", DOCKER_SERVER_ADDRESS,
		"-u", DOCKER_USERNAME,
		"-p", DOCKER_PASSWORD,
	)

	loginOut, err := loginCmd.CombinedOutput()
	if err != nil {
		ds.logger.Error("Docker login failed",
			zap.String("output", string(loginOut)),
			zap.Error(err))
		return err
	}
	ds.logger.Info("Docker login successful", zap.String("execution level", "ImagePush"))

	// Step 2: Push the image
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
