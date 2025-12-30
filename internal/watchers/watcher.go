package watchers

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
)

type Watcher struct{
	watcher *fsnotify.Watcher
}

func NewWatcher(watcher *fsnotify.Watcher) *Watcher {
	return &Watcher{watcher: watcher}
}

// WatchPeriodically adds sub directories peridically to watch, with the help
// of fsnotify which maintains a directory map rather than slice.
func (w *Watcher) WatchPeriodically(directory string, intervalSeconds int64) {
	done := make(chan struct{})
	go func() {
		done <- struct{}{}
	}()

	ticker := time.NewTicker(time.Duration(intervalSeconds) * time.Second)
	defer ticker.Stop()

	for ; ; <-ticker.C {
		<-done

		if err := filepath.Walk(directory, w.watchDir); err != nil {
			fmt.Println(err)
		}

		go func() {
			done <- struct{}{}
		}()
	}
}

// watchDir gets run as a walk func, searching for directories to add watchers to
func (w *Watcher) watchDir(path string, fi os.FileInfo, err error) error {
	// since fsnotify can watch all the files in a directory, watchers only need
	// to be added to each nested directory
	if fi.Mode().IsDir() {
		return w.watcher.Add(path)
	}

	return nil
}
