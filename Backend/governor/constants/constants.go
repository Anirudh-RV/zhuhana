package constants

import "os"

var UASAM_ENDPOINT = os.Getenv("UASAM_URL")
var FORGE_ENDPOINT = os.Getenv("FORGE_URL")
var MICROSERVICE_LOGIN_ENDPOINT = UASAM_ENDPOINT + "/v1/microservice/login/"
var API_AUTHENTICATION_ENDPOINT = UASAM_ENDPOINT + "/v1/microservice/authenticate/"
var USER_AUTHENTICATION_ENDPOINT = UASAM_ENDPOINT + "/v1/user/authenticate/"
var FORGE_BUILD_PYTHON_USER_ALGORITHM = FORGE_ENDPOINT + "/v1/algorithm/python/build/"
