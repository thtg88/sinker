package handlers

import (
	"context"
	"fmt"

	"github.com/fsnotify/fsnotify"

	"github.com/thtg88/sinker/internal/uploaders"
	"github.com/thtg88/sinker/pkg/sinker"
)

type FSEventHandler struct {
	fileUploader	uploaders.FileUploader
	sinkerAPI			sinker.API
}

func NewFSEventHandler(fileUploader uploaders.FileUploader, sinkerAPI sinker.API) *FSEventHandler {
	return &FSEventHandler{
		fileUploader: fileUploader,
		sinkerAPI: sinkerAPI,
	}
}

// HandleFsEvent handles a file system event, uploading a file to S3,
// and updates the state backend
func (h *FSEventHandler) Handle(event fsnotify.Event, sinkerAPIDeviceID string) {
	var err error

	// Skip CHMOD event as macOS sends 2 for every WRITE event (before and after)
	if event.Op.String() == "CHMOD" {
		return
	}

	ctx := context.Background()

	// TODO: replace with logger
	fmt.Printf("EVENT! %#v\n", event.String())

	if event.Op.String() == "CREATE" {
		err = h.fileUploader.UploadFile(ctx, event.Name)
	}

	if event.Op.String() == "REMOVE" {
		err = h.fileUploader.RemoveFile(ctx, event.Name)
	}

	if event.Op.String() == "WRITE" {
		err = h.fileUploader.UploadFile(ctx, event.Name)
	}

	if err != nil {
		// TODO: replace with logger
		fmt.Println("ERROR", err, event.Name)
		return
	}

	_, err = h.sinkerAPI.UpdateState(event, sinkerAPIDeviceID)
	if err != nil {
		// TODO: replace with logger
		fmt.Println("ERROR", err, event.Name)
	}

	// TODO: replace with logger
	fmt.Println("file updated", event.Name)
}
