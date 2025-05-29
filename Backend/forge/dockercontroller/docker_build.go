package dockercontroller

import (
	"fmt"
	"os"
	"os/exec"

	"go.uber.org/zap"
)

func (ds *DockerService) BuildImage(dockerfileDir, dockerImageName string) {
	fullImageName := fmt.Sprintf("%s/%s", DOCKER_USERNAME, dockerImageName)

	ds.logger.Info("Starting build image",
		zap.String("dockerfileDir", dockerfileDir),
		zap.String("imageName", fullImageName),
		zap.String("execution level", "BuildImage"),
	)

	// Prepare the docker build command
	cmd := exec.Command("docker", "build", "-t", fullImageName, dockerfileDir)

	// Set command output to stdout and stderr
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	if err := cmd.Run(); err != nil {
		ds.logger.Error("Docker build failed",
			zap.String("execution level", "BuildImage"),
			zap.Error(err),
		)
		return
	}

	ds.logger.Info("Docker image built successfully",
		zap.String("execution level", "BuildImage"),
	)
}
