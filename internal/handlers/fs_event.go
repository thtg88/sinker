package handlers

import (
	"fmt"

	"github.com/fsnotify/fsnotify"

	"github.com/thtg88/sinker/pkg/sinker"
	"github.com/thtg88/sinker/uploaders"
)

type FSEventHandler struct {
	sinkerAPI sinker.API
}

func NewFSEventHandler(sinkerAPI sinker.API) *FSEventHandler {
	return &FSEventHandler{sinkerAPI: sinkerAPI}
}

// HandleFsEvent handles a file system event, uploading a file to S3,
// and updates the state backend
func (h *FSEventHandler) Handle(event fsnotify.Event, sinkerAPIDeviceID string) {
	var err error

	// Skip CHMOD event as macOS sends 2 for every WRITE event (before and after)
	if event.Op.String() == "CHMOD" {
		return
	}

	fmt.Printf("EVENT! %#v\n", event.String())

	if event.Op.String() == "CREATE" {
		_, err = uploaders.UploadFile(event.Name)
	}

	if event.Op.String() == "REMOVE" {
		_, err = uploaders.RemoveFile(event.Name)
	}

	if event.Op.String() == "WRITE" {
		_, err = uploaders.UploadFile(event.Name)
	}

	if err != nil {
		fmt.Println("ERROR", err, event.Name)
		return
	}

	_, err = h.sinkerAPI.UpdateState(event, sinkerAPIDeviceID)
	if err != nil {
		fmt.Println("ERROR", err, event.Name)
	}

	fmt.Println("file updated", event.Name)
}
