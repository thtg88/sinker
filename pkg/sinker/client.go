package sinker

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/fsnotify/fsnotify"
)

type API interface {
	RegisterDevice() (string,error)
	UpdateState(event fsnotify.Event, sinkerAPIDeviceID string) ([]byte, error)
}

type APIClient struct{
	storeEventPath string
}

func NewAPIClient(storeEventPath string) *APIClient {
	return &APIClient{storeEventPath: storeEventPath}
}

func sinkerApiRequest(method string, uri string, requestBody []byte, sinkerAPIDeviceID string) ([]byte, error) {
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
