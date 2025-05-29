package dockercontroller

import "os"

var DOCKER_USERNAME = os.Getenv("DOCKER_USERNAME")
var DOCKER_PASSWORD = os.Getenv("DOCKER_PASSWORD")
var DOCKER_SERVER_ADDRESS = os.Getenv("DOCKER_SERVER_ADDRESS")
var DJANGO_TEMPLATE_PATH = "python-template"
