// pkg/thd/text_lognotify.go
package thd

import (
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

// Configuration
const (
	logDirectory   = "C:/Users/Public/GCS/Logging/"
	logFilePattern = "*.log"
	MaxWorkers     = 10
	logWriteEvent  = fsnotify.Write
	logCreateEvent = fsnotify.Create
	logRemoveEvent = fsnotify.Remove
)

var logger = log.New(os.Stdout, "LOG: ", log.Ldate|log.Ltime)

type LogFileInfo struct {
	Pattern string
	Path    string
	Name    string
}

// StartLogWatcher initializes a log file notifier to monitor changes in log files.
// It checks whether the log file exists and, if not, searches for another one that matches the specified pattern.
// Whenever a log file is found or changed, the information is sent to the provided channel fileInfoCh.
func StartLogWatcher(fileInfoCh chan<- LogFileInfo) {
	latestLogFile, err := getLatestLogFile()
	if err != nil {
		log.Println("Error:", err)
		return
	}

	if latestLogFile == "" {
		log.Println("No log files found matching the pattern:", logFilePattern)
	} else {
		log.Println("Latest log file found:", latestLogFile)
	}

	// Send the log file information to the channel
	fileInfo := LogFileInfo{
		Pattern: logFilePattern,
		Path:    logDirectory,
		Name:    latestLogFile[len(logDirectory):],
	}
	fileInfoCh <- fileInfo

	// Start watching the log file for changes
	go watchLogFile(fileInfoCh)
}

func getLatestLogFile() (string, error) {
	files, err := filepath.Glob(filepath.Join(logDirectory, logFilePattern))
	if err != nil {
		return "", err
	}
	if len(files) == 0 {
		return "", nil
	}

	latestLogFile := files[0]
	for _, file := range files {
		if file > latestLogFile {
			latestLogFile = file
		}
	}

	return latestLogFile, nil
}

func watchLogFile(fileInfoCh chan<- LogFileInfo) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	if err := watcher.Add(logDirectory); err != nil {
		return err
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}

			log.Printf("Event: %s, Op: %s", event.Name, event.Op)

			if event.Op&fsnotify.Write == fsnotify.Write {
				go processEvent(event, fileInfoCh)
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			log.Printf("Error: %v", err)
		}
	}
}

func processEvent(event fsnotify.Event, fileInfoCh chan<- LogFileInfo) {
	log.Printf("Processing event: %s, Op: %s", event.Name, event.Op)

	switch event.Op {
	case logWriteEvent, logCreateEvent:
		latestLogFile, err := getLatestLogFile()
		if err != nil {
			log.Println("Error:", err)
			return
		}
		fileInfo := LogFileInfo{
			Pattern: logFilePattern,
			Path:    logDirectory,
			Name:    latestLogFile[len(logDirectory):],
		}
		fileInfoCh <- fileInfo
	case logRemoveEvent:
		// Handle file removal event if needed
	}
}
