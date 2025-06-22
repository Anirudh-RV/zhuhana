package dockercontroller

import (
	"context"
	"fmt"
	"forge/logger"
	"os"
	"os/exec"

	"go.uber.org/zap"
)

type DockerService struct {
	logger *logger.Logger
	ctx    *context.Context
}

func NewDockerService(logger *logger.Logger) *DockerService {
	ctx := context.Background()

	// Login to Docker registry
	logger.Info("Logging in to Docker registry...",
		zap.String("registry", DOCKER_SERVER_ADDRESS),
		zap.String("execution level", "ImagePush"))
	loginCmd := exec.Command("docker", "login", DOCKER_SERVER_ADDRESS,
		"-u", DOCKER_USERNAME,
		"-p", DOCKER_PASSWORD,
	)

	loginOutput, err := loginCmd.CombinedOutput()
	if err != nil {
		logger.Fatal("Docker login failed",
			zap.String("output", string(loginOutput)),
			zap.Error(err))
	}
	logger.Info("Docker login successful", zap.String("execution level", "ImagePush"))

	return &DockerService{
		logger: logger,
		ctx:    &ctx,
	}
}

func (ds *DockerService) BuildUserAlgorithm(userID, scriptID, scriptURLPath string) error {
	destinationFolder := fmt.Sprintf("user-algorithm-%s-%s", userID, scriptID)
	destinationFilePath := fmt.Sprintf("%s/algorithm/zhuhana_algorithm.py", destinationFolder)
	dockerImageName := fmt.Sprintf("user-algorithm-%s-%s", userID, scriptID)
	defer os.RemoveAll(destinationFolder)
	defer ds.RemoveImage(dockerImageName)

	err := ds.CopyTemplate(destinationFolder)
	if err != nil {
		go ds.logger.Warning("copy template issue", zap.String("execution level", "BuildUserAlgorithm"), zap.String("Error", err.Error()))
		return err
	}
	err = ds.InsertUserScript(scriptURLPath, destinationFilePath)
	if err != nil {
		go ds.logger.Warning("insert user script issue", zap.String("execution level", "BuildUserAlgorithm"), zap.String("Error", err.Error()))
		return err
	}
	err = ds.BuildImage(destinationFolder, dockerImageName)
	if err != nil {
		go ds.logger.Warning("error while building image", zap.String("execution level", "BuildUserAlgorithm"), zap.String("Error", err.Error()))
		return err
	}
	ds.ImagePush(dockerImageName)

	return nil
}
