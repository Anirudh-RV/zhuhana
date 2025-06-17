package kubernetescontroller

import "os"

var DOCKER_USERNAME = os.Getenv("DOCKER_USERNAME")
var DOCKER_PASSWORD = os.Getenv("DOCKER_PASSWORD")
var DOCKER_REPOSITORY = os.Getenv("DOCKER_REPOSITORY")
var DOCKER_SERVER_ADDRESS = os.Getenv("DOCKER_SERVER_ADDRESS")

const KubeConfigPath = "/root/.kube/config"
const HostMinikubePath = "/Users/anirudhrv/.minikube"
const ContainerMinikubePath = "/root/.minikube"
const DefaultMinikubeServer = "https://127.0.0.1:62889"
const MinikubeIP = "192.168.49.2"
