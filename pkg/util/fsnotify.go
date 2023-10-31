// pkg/util/fsnotify.go
package util

import (
	"fmt"
	"log"
	"strings"
	"time"

	"testcode/pkg/thd"

	"github.com/fsnotify/fsnotify"
)

const (
	failedString            = "FAILED"
	changeInkString         = "change selected Ink to"
	initProgrammingString   = "Init Programming"
	verifiedString          = "verified"
	maxConcurrentInsertions = 5
	throttleDuration        = 500 * time.Millisecond
)

func FsnotifyStart(jobFull string, counterCallback func(string)) {
	// Initialize some constants and variables
	dirname := config.DirName
	joborder := jobFull
	initialStartup := false

	// Use a semaphore to limit concurrent insertions
	sem := make(chan struct{}, maxConcurrentInsertions)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println("Error creating watcher:", err)
		return
	}
	defer watcher.Close()

	done := make(chan bool)

	err = watcher.Add(dirname)
	if err != nil {
		log.Println("Error adding watcher:", err)
		return
	}

	log.Println("監視開始", dirname)

	// Start a separate goroutine to handle file system events
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					close(done)
					return
				}
				handleEvent(event, &initialStartup, joborder, counterCallback)
			case err, ok := <-watcher.Errors:
				if !ok {
					close(done)
					return
				}
				log.Println("Error:", err)
			}
		}
	}()

	// Start a separate goroutine to handle thd.Insert operations
	go func() {
		fmt.Println(joborder)
		for contents := range contentCh {
			if initialStartup && !shouldSkip(contents) {
				// Acquire a semaphore
				sem <- struct{}{}
				go func(contents string) {
					defer func() {
						// Release the semaphore
						<-sem
					}()
					if err := insertData(contents, joborder, counterCallback); err != nil {
						log.Println("Error inserting data:", err)
					}
					// Clear Contents and joborder
					SetContents("")
				}(contents)
			}
		}
	}()

	<-done
}

// Ensure that the contents are not empty and then pass them to SetContents.
func handleEvent(event fsnotify.Event, initialStartup *bool, joborder string, counterCallback func(string)) {
	if event.Op&fsnotify.Write == fsnotify.Write {
		*initialStartup = true
		contents, contentsReady := getContents()
		if !contentsReady {
			// Contents not ready, wait until it becomes available
			contents, contentsReady = waitForContentsReady()
		}
		SetContents(contents)
	}
}

// Condition for thd.Insert: Skip the insertion session if any of the following keywords are encountered.
func shouldSkip(contents string) bool {
	return strings.Contains(contents, failedString) ||
		strings.Contains(contents, initProgrammingString) ||
		strings.Contains(contents, changeInkString) ||
		!strings.Contains(contents, verifiedString)
}

// Function for inserting the log's content into the database
func insertData(contents, joborder string, counterCallback func(string)) error {
	//log.Printf("fsnotify.go: %s", contents)
	sig, err := thd.Insert(
		config.DBHost,
		config.DBPort,
		config.DBUser,
		config.DBPassword,
		config.DBName,
		config.PCName,
		joborder,
		contents,
	)
	if err != nil {
		return err
	}
	// Call the completion callback
	counterCallback(sig)
	return nil
}

func SetContents(value string) {
	contentCh <- value
}

func waitForContentsReady() (string, bool) {
	contents := <-contentCh
	return contents, true
}

func getContents() (string, bool) {
	contentMutex.Lock()
	defer contentMutex.Unlock()
	return content, contentReady
}
