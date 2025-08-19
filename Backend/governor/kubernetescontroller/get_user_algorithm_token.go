package kubernetescontroller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"governor/constants"
	"governor/user/algorithm/models"
	"io"
	"net/http"
)

func (ks *KubernetesService) GetUserAlgorithmToken(userAlgorithmID string) (string, error) {

	bodyData := map[string]string{
		"userAlgorithmID": userAlgorithmID,
	}
	bodyBytes, err := json.Marshal(bodyData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal body: %w", err)
	}

	// Create request
	req, err := http.NewRequest("POST", constants.MICROSERVICE_USER_ALGORITHM_LOGIN_ENDPOINT, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("AUTH_TOKEN", ks.microserviceAuthenticator.ALL_SERVICE_JWT_TOKENS[ks.microserviceAuthenticator.UASAM_SERVICE_NAME])
	req.Header.Set("ORIGIN_SERVICE", ks.microserviceAuthenticator.ORIGIN_SERVICE)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var loginResp models.UserAlgorithmLoginResponse
	if err := json.Unmarshal(respBody, &loginResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if loginResp.Status != 1 {
		return "", fmt.Errorf("login failed : %s", loginResp.StatusDescription)
	}

	return loginResp.AccessToken, nil
}
