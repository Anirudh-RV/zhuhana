package dockercontroller

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

func (ds *DockerService) InsertUserScript(url, destinationPath string) error {
	// Send GET request
	resp, err := http.Get(url)
	if err != nil {
		return errors.New(fmt.Sprint("failed to fetch url: %w", err))
	}
	defer resp.Body.Close()

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		return errors.New("bad status")
	}

	// Create destination file
	out, err := os.Create(destinationPath)
	if err != nil {
		return errors.New(fmt.Sprint("failed to create file: %w", err))
	}
	defer out.Close()

	// Copy response body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return errors.New(fmt.Sprint("failed to write file: %w", err))
	}

	return nil
}

// copyFile copies a single file from src to dst
func (ds *DockerService) copyFile(src, dst string, info fs.FileInfo) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// Copy content
	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	// Copy permissions
	return os.Chmod(dst, info.Mode())
}

// copyDir recursively copies a directory tree from src to dst
func (ds *DockerService) copyDir(src string, dst string) error {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// Create destination directory
	err = os.MkdirAll(dst, srcInfo.Mode())
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		entryInfo, err := entry.Info()
		if err != nil {
			return err
		}

		if entry.IsDir() {
			// Recursively copy sub-directory
			if err := ds.copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			// Copy file
			if err := ds.copyFile(srcPath, dstPath, entryInfo); err != nil {
				return err
			}
		}
	}
	return nil
}

func (ds *DockerService) CopyTemplate(destinationFolder string) error {
	err := ds.copyDir(DJANGO_TEMPLATE_PATH, destinationFolder)
	if err != nil {
		go ds.logger.Warning("copydir error", zap.String("execution level", "CopyTemplate"), zap.String("Error", err.Error()))
		return err
	}
	return nil
}
