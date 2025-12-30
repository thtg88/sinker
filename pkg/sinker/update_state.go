package sinker

import (
	"encoding/json"

	"github.com/fsnotify/fsnotify"
	"github.com/thtg88/sinker/uploaders"
)

// UpdateState updates the state backend
func (c *APIClient) UpdateState(event fsnotify.Event, sinkerAPIDeviceID string) ([]byte, error) {
	values := map[string]string{
		"path": uploaders.RelativePath(event.Name),
		"type": event.Op.String(),
	}
	jsonValue, _ := json.Marshal(values)

	return sinkerApiRequest("POST", c.storeEventPath, jsonValue, sinkerAPIDeviceID)
}
