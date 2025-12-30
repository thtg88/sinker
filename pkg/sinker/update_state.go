package sinker

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type updateStateRequest struct {
	Path string `json:"path"`
	Type string `json:"type"`
}

// UpdateState updates the state backend
func (c *APIClient) UpdateState(relativePath string, operation string, sinkerAPIDeviceID string) ([]byte, error) {
	request := updateStateRequest{
		Path: relativePath,
		Type: operation,
	}
	jsonValue, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("json marshal: %v", err)
	}

	return c.sinkerApiRequest(http.MethodPost, c.config.StoreEventPath, jsonValue, sinkerAPIDeviceID)
}
