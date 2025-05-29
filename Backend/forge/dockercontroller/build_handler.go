package dockercontroller

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"go.uber.org/zap"
)

func (ds *DockerService) BuildImage(dockerfileDir, dockerImageName string) error {
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
		return err
	}

	ds.logger.Info("Docker image built successfully",
		zap.String("execution level", "BuildImage"),
	)

	return nil
}

func (ds *DockerService) getImageID(imageName string) (string, error) {
	cmd := exec.Command("docker", "inspect", "--format={{.Id}}", imageName)
	ds.logger.Info("image id retreival: ", zap.String("Running command: %s\n", strings.Join(cmd.Args, " ")))

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get image ID: %w, output: %s", err, string(output))
	}
	imageID := strings.TrimSpace(string(output))
	return imageID, nil
}

func (ds *DockerService) RemoveImage(dockerImageName string) error {
	fullImageName := fmt.Sprintf("%s/%s", DOCKER_USERNAME, dockerImageName)
	imageID, err := ds.getImageID(fullImageName)
	if err != nil {
		ds.logger.Error("docker image removal failed",
			zap.String("execution level", "RemoveImage"),
			zap.Error(err),
		)
		return err
	}

	ds.logger.Info("Starting image untag",
		zap.String("imageName", fullImageName),
		zap.String("execution level", "RemoveImage"),
	)

	// Step 1: Untag the image
	cmd := exec.Command("docker", "rmi", "-f", fullImageName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		ds.logger.Error("docker image removal failed",
			zap.String("execution level", "RemoveImage"),
			zap.Error(err),
		)
		return err
	}

	ds.logger.Info("starting image removal",
		zap.String("imageName", imageID),
		zap.String("execution level", "RemoveImage"),
	)
	// Step 2: Remove the image
	cmd = exec.Command("docker", "rmi", "-f", imageID)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		ds.logger.Error("docker image removal failed",
			zap.String("execution level", "RemoveImage"),
			zap.Error(err),
		)
		return err
	}

	ds.logger.Info("image removed successfully",
		zap.String("execution level", "RemoveImage"),
	)

	return nil
}
