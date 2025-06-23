package constants

import "os"

var UASAM_ENDPOINT = os.Getenv("UASAM_URL")
var ALGONEXUS_URL = os.Getenv("ALGONEXUS_URL")
var MICROSERVICE_LOGIN_ENDPOINT = UASAM_ENDPOINT + "/v1/microservice/login/"
var API_AUTHENTICATION_ENDPOINT = UASAM_ENDPOINT + "/v1/microservice/authenticate/"
var MICROSERVICE_USER_ALGORITHM_AUTHENTICATE_ENDPOINT = UASAM_ENDPOINT + "/v1/microservice/user-algorithm/authenticate/"
