package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/joho/godotenv"

	"github.com/thtg88/sinker/internal/handlers"
	"github.com/thtg88/sinker/internal/watchers"
	"github.com/thtg88/sinker/pkg/sinker"
)

var watcher *fsnotify.Watcher
var sinkerAPIDeviceID string

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(errors.New("could not load .env file"))
	}

	sinkerAPIDeviceID, err = sinker.RegisterDevice()
	if err != nil {
		panic(err)
	}

	// TODO sync before watching

	fsNotifyWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	defer fsNotifyWatcher.Close()

	watcher := watchers.NewWatcher(fsNotifyWatcher)

	go watcher.WatchPeriodically(os.Getenv("SINKER_BASE_PATH"), 5)

	done := make(chan bool)

	go func() {
		for {
			select {
			case event := <-fsNotifyWatcher.Events:
				handlers.HandleFsEvent(event, sinkerAPIDeviceID)

			case err := <-fsNotifyWatcher.Errors:
				fmt.Println("ERROR", err)
			}
		}
	}()

	<-done
}
