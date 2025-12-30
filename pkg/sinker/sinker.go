package sinker

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/thtg88/sinker/uploaders"
)

// UpdateState updates the state backend
func UpdateState(event fsnotify.Event,sinkerAPIDeviceID string) ([]byte, error) {
	values := map[string]string{
		"path": uploaders.RelativePath(event.Name),
		"type": event.Op.String(),
	}
	jsonValue, _ := json.Marshal(values)

	return SinkerApiRequest("POST", os.Getenv("SINKER_API_STORE_EVENT_PATH"), jsonValue, sinkerAPIDeviceID)
}

func RegisterDevice() (string, error) {
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

	body, err := SinkerApiRequest("POST", os.Getenv("SINKER_API_STORE_DEVICE_PATH"), jsonValue, "")
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

func SinkerApiRequest(method string, uri string, requestBody []byte, sinkerAPIDeviceID string) ([]byte, error) {
	url := fmt.Sprint(os.Getenv("SINKER_API_BASE_URL"), uri)
	fmt.Println("URL:>", url)

	req, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set(os.Getenv("SINKER_API_KEY_HEADER_NAME"), os.Getenv("SINKER_API_KEY_HEADER_VALUE"))
	req.Header.Set(os.Getenv("SINKER_API_USER_ID_HEADER_NAME"), os.Getenv("SINKER_API_USER_ID_HEADER_VALUE"))
	req.Header.Set(os.Getenv("SINKER_API_DEVICE_ID_HEADER_NAME"), sinkerAPIDeviceID)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// fmt.Println("request Headers:", req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
	return body, nil
}
