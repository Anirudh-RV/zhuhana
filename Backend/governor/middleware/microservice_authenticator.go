package middleware

import (
	"encoding/json"
	"fmt"
	"governor/constants"
	"governor/logger"
	"io"
	"net/http"
	"os"

	"go.uber.org/zap"
)

type MicroSeviceAuthenticator struct {
	logger                 *logger.Logger
	ORIGIN_SERVICE         string
	ALL_API_KEYS           map[string]string
	ALL_SERVICE_JWT_TOKENS map[string]string
	FORGE_SERVICE_NAME     string
}

type MicroServiceLoginResponse struct {
	Status            int    `json:"status"`
	StatusDescription string `json:"statusDescription"`
	AccessToken       string `json:"accessToken,omitempty"`
}

func NewMicroSeviceAuthenticator(logger *logger.Logger) *MicroSeviceAuthenticator {
	ORIGIN_SERVICE := os.Getenv("ORIGIN_SERVICE")
	FORGE_API_KEY := os.Getenv("FORGE_API_KEY")
	ALL_SERVICE_JWT_TOKENS := map[string]string{}

	ALL_API_KEYS := map[string]string{
		FORGE_API_KEY: "forge",
	}

	return &MicroSeviceAuthenticator{
		logger:                 logger,
		ORIGIN_SERVICE:         ORIGIN_SERVICE,
		ALL_API_KEYS:           ALL_API_KEYS,
		FORGE_SERVICE_NAME:     "forge",
		ALL_SERVICE_JWT_TOKENS: ALL_SERVICE_JWT_TOKENS,
	}
}

func (msa *MicroSeviceAuthenticator) GetAllServiceTokens() error {
	for apiKey, serviceName := range msa.ALL_API_KEYS {
		go msa.logger.Info(fmt.Sprintf("getting jwt token service %s", serviceName), zap.String("execution level", "GetAllServiceTokens"))
		serviceJWTToken, err := msa.GetServiceTokens(apiKey)
		if err != nil {
			go msa.logger.Fatal(fmt.Sprintf("could not get jwt token for %s", serviceName), zap.String("execution level", "GetAllServiceTokens"), zap.String("Error", err.Error()))
		}
		msa.ALL_SERVICE_JWT_TOKENS[serviceName] = serviceJWTToken
	}

	return nil
}

func (msa *MicroSeviceAuthenticator) GetServiceTokens(apiKey string) (string, error) {
	req, err := http.NewRequest("POST", constants.MICROSERVICE_LOGIN_ENDPOINT, nil)
	if err != nil {
		return "", err
	}

	// Add headers
	req.Header.Add("API_KEY", apiKey)
	req.Header.Add("ORIGIN_SERVICE", msa.ORIGIN_SERVICE)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("unexpected status code: %d\nresponse body: %s", resp.StatusCode, string(bodyBytes))
	}

	// Decode the response body into the struct
	var response MicroServiceLoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}

	if response.Status != 1 {
		return "", fmt.Errorf("error getting jwt token %s", response.StatusDescription)
	}

	return response.AccessToken, nil
}
