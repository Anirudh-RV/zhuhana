package commonutils

import "mime/multipart"

func UploadFileToCloudStorage(userID, scriptID string, script multipart.File) (string, error) {
	// TODO
	// ONCE CLOUD PROVIDER IS FINALIZED
	return "https://raw.githubusercontent.com/Anirudh-RV/UploadToS3/refs/heads/main/runner.py", nil
}

func GetPresignedURL(userID, scriptID string) (string, error) {
	// TODO
	// ONCE CLOUD PROVIDER IS FINALIZED
	return "https://raw.githubusercontent.com/Anirudh-RV/UploadToS3/refs/heads/main/runner.py", nil
}
