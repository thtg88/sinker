package handlers

import (
	"context"
	"log"

	"github.com/fsnotify/fsnotify"

	"github.com/thtg88/sinker/internal/uploaders"
	"github.com/thtg88/sinker/pkg/sinker"
)

type FSEventHandler struct {
	fileUploader	uploaders.FileUploader
	logger				*log.Logger
	sinkerAPI			sinker.API
}

func NewFSEventHandler(fileUploader uploaders.FileUploader, sinkerAPI sinker.API, logger *log.Logger) *FSEventHandler {
	return &FSEventHandler{
		fileUploader: fileUploader,
		logger:				logger,
		sinkerAPI:		sinkerAPI,
	}
}

// HandleFsEvent handles a file system event, uploading a file to S3,
// and updates the state backend
func (h *FSEventHandler) Handle(event fsnotify.Event, sinkerAPIDeviceID string) {
	var err error

	ctx := context.Background()

	opString := event.Op.String()

	h.logger.Printf("event: %s", event.String())

	if event.Has(fsnotify.Create) || event.Has(fsnotify.Write) {
		err = h.fileUploader.UploadFile(ctx, event.Name)
	} else if event.Has(fsnotify.Remove) {
		err = h.fileUploader.RemoveFile(ctx, event.Name)
	} else if event.Has(fsnotify.Chmod) {
		// Skip CHMOD event as macOS sends 2 for every WRITE event (before and after)
		return
	}

	if err != nil {
		h.logger.Printf("[ERROR] fileuploader operation: %s event name: %s: %v", opString, event.Name, err)
		return
	}

	_, err = h.sinkerAPI.UpdateState(h.fileUploader.RelativePath(event.Name), opString, sinkerAPIDeviceID)
	if err != nil {
		h.logger.Printf("[ERROR] sinkerapi updatestate operation: %s event name: %s: %v", opString, event.Name, err)
	}

	h.logger.Printf("file updated operation: %s event name: %s", opString, event.Name)
}
