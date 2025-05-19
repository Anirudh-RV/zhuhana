package constants

import "os"

var UASAM_ENDPOINT = os.Getenv("UASAM_URL")
var API_AUTHENTICATION_ENDPOINT = UASAM_ENDPOINT + "/v1/microservice/authenticate/"
