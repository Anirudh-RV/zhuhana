package services

import (
	"forge/dockercontroller"
	"forge/logger"

	"go.uber.org/zap"
)

type PythonBuilderService struct {
	logger        *logger.Logger
	dockerService *dockercontroller.DockerService
}

func NewPythonBuilderService(logger *logger.Logger, dockerService *dockercontroller.DockerService) *PythonBuilderService {
	return &PythonBuilderService{
		logger:        logger,
		dockerService: dockerService,
	}
}

func (pbs *PythonBuilderService) BuildAlgorithmHandler(userID, scriptID, scriptURL string) error {

	err := pbs.dockerService.BuildUserAlgorithm(userID, scriptID, scriptURL)
	if err != nil {
		go pbs.logger.Warning("error while building and pushing user algorithm", zap.String("execution level", "BuildAlgorithmHandler"), zap.String("Error", err.Error()))
		return err
	}
	return nil
}
