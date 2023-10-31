package log

import (
	"fmt"
	"log"
	"sync"
	"time"

	thd "testcode/pkg/thd"
	util "testcode/pkg/util"
)

var (
	logPathMutex      sync.Mutex
	logPath           string
	lastProcessedFile string
)

const logWatcherDelay = 50 * time.Millisecond

func PickLog(filePattern string, fileInfo thd.LogFileInfo, findStr string) (string, error) {
	contents, err := thd.GetLog(fileInfo.Name, findStr, fileInfo.Path)
	if err != nil {
		return "", fmt.Errorf("error getting log for %s: %w", fileInfo.Name, err)
	}
	return contents, nil
}

func LogWatcherLoop(filePattern string, findStr string) {
	var wg sync.WaitGroup
	// Create a channel for receiving log file information
	fileInfoCh := make(chan thd.LogFileInfo, 300)

	// Start the log watcher with the channel
	go thd.StartLogWatcher(fileInfoCh)
	wg.Add(1)
	// Loop continuously to process log files
	var logContentBuffer string
	var lastLogUpdateTime time.Time
	// Use a channel to signal the processing routine
	processCh := make(chan struct{})

	go func() {
		for {
			fileInfo := <-fileInfoCh
			logContent, err := PickLog(filePattern, fileInfo, findStr)
			if err != nil {
				log.Printf("error getting log for %s: %v", fileInfo.Name, err)
				// Consider adding an error handling strategy here
				continue
			}
			//log.Printf("log.go: %s", logContent)
			logPathMutex.Lock()
			logPath = fileInfo.Path + fileInfo.Name
			logPathMutex.Unlock()

			// Store the log content in the buffer
			logContentBuffer = logContent
			lastLogUpdateTime = time.Now()

			// Signal the processing routine
			processCh <- struct{}{}

			time.Sleep(10 * time.Millisecond)
		}
	}()

	// Process the buffered content when signaled
	go func() {
		for range processCh {
			// Check if there's new log content within the last interval
			if time.Since(lastLogUpdateTime) < logWatcherDelay {
				util.SetContents(logContentBuffer)
			}
		}
	}()

	wg.Wait()
}
