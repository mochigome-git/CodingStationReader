// pkg/util/log.go

package util

import (
	"os"
)

func SetLogPath(value string) {
	logPathMutex.Lock()
	defer logPathMutex.Unlock()
	logPath = value
}

func GetLogPath() (string, bool) {
	logPathMutex.Lock()
	defer logPathMutex.Unlock()
	return logPath, logPath != ""
}

func DeleteLogFileContents(filePath string) error {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString("")
	if err != nil {
		return err
	}

	return nil
}
