package sinker

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

func (c *APIClient) RegisterDevice() (string, error) {
	type Device struct {
		Uuid string `json:"uuid"`
		Name string `json:"name"`
	}
	type DeviceResponse struct {
		Data Device `json:"data"`
	}
	var response DeviceResponse

	values := map[string]string{"name": "Sinker macOS client"}
	jsonValue, _ := json.Marshal(values)

	body, err := sinkerApiRequest("POST", os.Getenv("SINKER_API_STORE_DEVICE_PATH"), jsonValue, "")
	if err != nil {
		return "", err
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return "", err
	}
	if response.Data.Uuid == "" {
		return "", errors.New(fmt.Sprint("Could not register device. ERROR ", string(body)))
	}

	return string(response.Data.Uuid), nil
}
