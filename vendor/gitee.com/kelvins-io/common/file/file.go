package file

import (
	"os"
	"path/filepath"
)

// get all files under path
func GetFileList(path string) ([]string, error) {
	var fileList []string
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		fileList = append(fileList, path)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return fileList, nil
}

// judge file is exist
func IsFileExists(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
