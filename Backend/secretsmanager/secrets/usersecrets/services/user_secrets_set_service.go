package services

import "go.uber.org/zap"

func (uss *UserSecretsService) SetUserSecret(userID, key, value string) error {

	err := uss.userSecretsRepository.CreateUserSecret(userID, key, value)
	if err != nil {
		go uss.logger.Warning("error setting user secret", zap.String("execution level", "SetUserSecrets"), zap.String("Error", err.Error()))
		return err
	}

	return nil
}
