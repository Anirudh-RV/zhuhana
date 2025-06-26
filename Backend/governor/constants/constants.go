package constants

import "os"

var UASAM_ENDPOINT = os.Getenv("UASAM_URL")
var FORGE_ENDPOINT = os.Getenv("FORGE_URL")
var USER_ALGORITHM_API_ENDPOINT = os.Getenv("USER_ALGORITHM_API_ENDPOINT") // Subject to change
var MICROSERVICE_LOGIN_ENDPOINT = UASAM_ENDPOINT + "/v1/microservice/login/"
var API_AUTHENTICATION_ENDPOINT = UASAM_ENDPOINT + "/v1/microservice/authenticate/"
var USER_AUTHENTICATION_ENDPOINT = UASAM_ENDPOINT + "/v1/user/authenticate/"
var FORGE_BUILD_PYTHON_USER_ALGORITHM = FORGE_ENDPOINT + "/v1/algorithm/python/build/"
var MICROSERVICE_USER_ALGORITHM_LOGIN_ENDPOINT = UASAM_ENDPOINT + "/v1/microservice/user-algorithm/login/"
var MICROSERVICE_USER_ALGORITHM_AUTHENTICATE_ENDPOINT = UASAM_ENDPOINT + "/v1/microservice/user-algorithm/authenticate/"
