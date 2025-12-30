package sinker

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const deviceName string = "Sinker macOS client"

type deviceRequest struct {
	Name string `json:"name"`
}

type deviceResponse struct {
	Data *device `json:"data"`
}

type device struct {
	Uuid string `json:"uuid"`
	Name string `json:"name"`
}

func (c *APIClient) RegisterDevice() (string, error) {
	request := deviceRequest{Name: deviceName}
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("json marshal: %v", err)
	}

	responseBytes, err := c.sinkerApiRequest(http.MethodPost, c.config.StoreDevicePath, requestBytes, "")
	if err != nil {
		return "", fmt.Errorf("sinker api request: %v", err)
	}

	var response deviceResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return "", fmt.Errorf("json unmarshal: %v", err)
	}
	if response.Data.Uuid == "" {
		return "", fmt.Errorf("empty response data uuid: %s", string(responseBytes))
	}

	return response.Data.Uuid, nil
}
