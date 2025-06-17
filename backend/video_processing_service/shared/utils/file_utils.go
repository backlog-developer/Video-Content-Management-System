// / File: shared/utils/file_utils.go
package utils

import (
	"os"
)

// EnsureDir makes sure the directory exists
func EnsureDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, os.ModePerm)
	}
	return nil
}
