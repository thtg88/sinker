package sinker

import (
	"encoding/json"
	"fmt"

	"github.com/fsnotify/fsnotify"

	"github.com/thtg88/sinker/internal/uploaders"
)

type updateStateRequest struct {
	Path string `json:"path"`
	Type string `json:"type"`
}

// UpdateState updates the state backend
func (c *APIClient) UpdateState(event fsnotify.Event, sinkerAPIDeviceID string) ([]byte, error) {
	request := updateStateRequest{
		Path: uploaders.RelativePath(event.Name),
		Type: event.Op.String(),
	}
	jsonValue, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("json marshal: %v", err)
	}

	return c.sinkerApiRequest("POST", c.config.StoreEventPath, jsonValue, sinkerAPIDeviceID)
}
