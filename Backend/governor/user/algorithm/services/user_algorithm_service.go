package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	commonutils "governor/commonutils"
	"governor/constants"
	"governor/logger"
	"governor/middleware"
	"governor/user/algorithm/models"
	"governor/user/algorithm/repositories"
	"io"
	"mime/multipart"
	"net/http"

	"go.uber.org/zap"
)

type UserAlgorithmService struct {
	logger                    *logger.Logger
	userAlgorthmRepository    *repositories.UserAlgorithmRepository
	microserviceAuthenticator *middleware.MicroSeviceAuthenticator
}

func NewUserAlgorithmService(logger *logger.Logger, userAlgorthmRepository *repositories.UserAlgorithmRepository, microserviceAuthenticator *middleware.MicroSeviceAuthenticator) *UserAlgorithmService {
	return &UserAlgorithmService{
		logger:                    logger,
		userAlgorthmRepository:    userAlgorthmRepository,
		microserviceAuthenticator: microserviceAuthenticator,
	}
}

func (uas *UserAlgorithmService) CreateUserAlgorithmHandler(userID, scriptName string, script multipart.File) (*models.UserAlgorithm, error) {
	// TODO
	userAlgorithm, err := uas.userAlgorthmRepository.CreateUserAlgorithm(userID, scriptName)
	if err != nil {
		go uas.logger.Error("could not create user algorithm entry", zap.String("execution level", "CreateUserAlgorithm"), zap.String("Error", err.Error()))
		return nil, err
	}

	go uas.logger.Info("spawning thread to upload, build and push user algorithm", zap.String("execution level", "CreateUserAlgorithm"))
	go uas.UploadUserAlgorithmScript(userID, userAlgorithm.ID.String(), script)
	return userAlgorithm, nil
}

func (uas *UserAlgorithmService) UploadUserAlgorithmScript(userID, scriptID string, script multipart.File) error {
	scriptURL, err := commonutils.UploadFileToCloudStorage(userID, scriptID, script)
	if err != nil {
		go uas.logger.Error("could not upload file to cloud storage", zap.String("execution level", "UploadUserAlgorithmScript"), zap.String("Error", err.Error()))
		return err
	}
	go uas.logger.Info(fmt.Sprintf("file uploaded to cloud storage: %s", scriptURL), zap.String("execution level", "UploadUserAlgorithmScript"))

	presignedScriptURL, err := commonutils.GetPresignedURL(userID, scriptID)
	if err != nil {
		go uas.logger.Error("could not get presigned url", zap.String("execution level", "UploadUserAlgorithmScript"), zap.String("Error", err.Error()))
		return err
	}
	go uas.logger.Info(fmt.Sprintf("presigned url created: %s", presignedScriptURL), zap.String("execution level", "UploadUserAlgorithmScript"))

	err = uas.userAlgorthmRepository.UpdateScriptURL(scriptID, scriptURL)
	if err != nil {
		go uas.logger.Error("could not update script url in the database", zap.String("execution level", "UploadUserAlgorithmScript"), zap.String("Error", err.Error()))
		return err
	}
	go uas.logger.Info(fmt.Sprintf("updated script url in the database: %s", scriptID), zap.String("execution level", "UploadUserAlgorithmScript"))

	err = uas.BuildAndPushUserAlgorithmScript(userID, scriptID, presignedScriptURL)
	if err != nil {
		go uas.logger.Error("could not build and push user algorithm script", zap.String("execution level", "UploadUserAlgorithmScript"), zap.String("Error", err.Error()))
		return err
	}

	go uas.logger.Info(fmt.Sprintf("built and pushed user algorithm: user-algorithm-%s-%s", userID, scriptID), zap.String("execution level", "UploadUserAlgorithmScript"))

	return nil
}

func (uas *UserAlgorithmService) BuildAndPushUserAlgorithmScript(userID, scriptID, scriptURL string) error {
	payload := models.PythonBuilderRequest{
		UserID:    userID,
		ScriptID:  scriptID,
		ScriptURL: scriptURL,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Create request
	req, err := http.NewRequest("POST", constants.FORGE_BUILD_PYTHON_USER_ALGORITHM, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("AUTH_TOKEN", uas.microserviceAuthenticator.ALL_SERVICE_JWT_TOKENS[uas.microserviceAuthenticator.FORGE_SERVICE_NAME])
	req.Header.Set("ORIGIN_SERVICE", uas.microserviceAuthenticator.ORIGIN_SERVICE)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d\nresponse body: %s", resp.StatusCode, string(bodyBytes))
	}

	// Decode JSON response
	var response models.PythonBuilderResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}
	if response.Status != 1 {
		return fmt.Errorf("failed to build and push container: %s", response.StatusDescription)
	}

	return nil
}
