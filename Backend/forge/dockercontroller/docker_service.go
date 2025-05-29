package dockercontroller

import (
	"context"
	"forge/logger"
)

type DockerService struct {
	logger *logger.Logger
	ctx    *context.Context
}

func NewDockerService(logger *logger.Logger) *DockerService {
	ctx := context.Background()

	return &DockerService{
		logger: logger,
		ctx:    &ctx,
	}
}
