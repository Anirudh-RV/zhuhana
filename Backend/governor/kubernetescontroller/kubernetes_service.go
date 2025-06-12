package kubernetescontroller

import (
	"database/sql"
	"fmt"
	"governor/logger"
	"os/exec"

	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type KubernetesService struct {
	logger    *logger.Logger
	db        *sql.DB
	clientSet *kubernetes.Clientset
	namespace string
}

func NewKubernetesService(logger *logger.Logger, db *sql.DB) *KubernetesService {
	// TODO: Login to Docker registry
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

	config, err := rest.InClusterConfig()
	if err != nil {
		logger.Fatal("Failed to get in-cluster config", zap.Error(err))
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	fmt.Printf("clientset connection completed")

	versionInfo, err := clientSet.Discovery().ServerVersion()
	if err != nil {
		fmt.Printf("unable to connect to cluster: %s", err.Error())
	}

	fmt.Printf("Connected to Kubernetes cluster version: %s\n", versionInfo.GitVersion)

	return &KubernetesService{
		logger:    logger,
		db:        db,
		clientSet: clientSet,
		namespace: "user-algorithm",
	}
}
