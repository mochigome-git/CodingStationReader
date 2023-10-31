package main

import (
	"embed"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"time"

	"github.com/zserge/lorca"

	_log "testcode/log"
	"testcode/pkg/util"
)

type Config struct {
	PCName  string
	DirName string
}

//go:embed www
var fs embed.FS

// Go types that are bound to the UI must be thread-safe because each binding
// is executed in its own goroutine. In this simple case, we may use atomic
// operations, but for more complex cases one should use proper synchronization.

var (
	jobFullMutex sync.Mutex
	// Create a channel to update jobFull
	jobFullChan = make(chan string, 1)
	// new channel to signal when the jobFull value has been updated
	jobFullUpdated = make(chan struct{})
	// job order
	jobFull         string
	fsnotifyStarted bool
)

func main() {
	// Initialize the configuration
	util.InitConfig()
	// Retrieve the config instance
	util.GetConfig()
	filePattern := time.Now().Format("2006_01_02") + "*.log"
	findStr := "read: "
	go _log.LogWatcherLoop(filePattern, findStr)

	args := []string{}
	if runtime.GOOS == "linux" {
		args = append(args, "--class=Lorca")
	}

	ui, err := lorca.New("", "", 780, 620, "--remote-allow-origins=*")
	if err != nil {
		log.Println("runtime error:", err)
	}
	defer ui.Close()

	c := &util.Counter{}

	ui.Bind("start", func() {
		log.Println("UI is ready")
	})

	ui.Bind("saveJobGo", func(jobOrder, jobMonth string) {
		jobFullMutex.Lock()
		defer jobFullMutex.Unlock()
		jobFull := jobMonth + jobOrder
		// Notify that jobFull has been updated
		jobFullChan <- jobFull
	})

	ui.Bind("counterAdd", c.Add)
	ui.Bind("counterValue", c.Count)
	ui.Bind("saveConfigGo", util.SetConfig)
	ui.Bind("loadConfigGo", util.LoadConfigGo)
	ui.Bind("resetCounter", c.Reset)
	ui.Bind("NotifyStart", func() {
		// Wait for the updated jobFull value
		updatedJobFull := <-jobFullChan
		if updatedJobFull != "" {
			go util.FsnotifyStart(updatedJobFull, func(sig string) {
				ui.Eval(`window.NotifyStartComplete("` + sig + `")`)
			})
			// Update the jobFull value and notify the monitoring goroutine
			jobFullMutex.Lock()
			jobFull = updatedJobFull
			jobFullMutex.Unlock()
			jobFullUpdated <- struct{}{}
			fsnotifyStarted = true
		}
	})

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Println("listen tcp error:", err)
	}
	defer ln.Close()
	go http.Serve(ln, http.FileServer(http.FS(fs)))
	ui.Load(fmt.Sprintf("http://%s/www", ln.Addr()))

	sigc := make(chan os.Signal)
	signal.Notify(sigc, os.Interrupt)
	go func() {
		<-sigc
		log.Println("exiting...")
		os.Exit(0)
	}()

	select {
	case <-sigc:
	case <-ui.Done():
		log.Println("exiting...")
		os.Exit(0)
	}
	log.Println("exiting...")
	os.Exit(0)
}
