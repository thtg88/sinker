package sinker

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/fsnotify/fsnotify"

	"github.com/thtg88/sinker/internal/config"
)

type API interface {
	RegisterDevice() (string,error)
	UpdateState(event fsnotify.Event, sinkerAPIDeviceID string) ([]byte, error)
}

type APIClient struct{
	httpClient	*http.Client
	config 			*config.SinkerAPI
}

func NewAPIClient(httpClient *http.Client, cfg *config.SinkerAPI) *APIClient {
	return &APIClient{
		config:			cfg,
		httpClient:	httpClient,
	}
}

func (c *APIClient) sinkerApiRequest(method string, uri string, requestBody []byte, sinkerAPIDeviceID string) ([]byte, error) {
	url := fmt.Sprint(c.config.BaseURL, uri)
	// TODO: replace with logger
	fmt.Println("URL:>", url)

	req, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("http newrequest: %v", err)
	}

	req.Header.Set(c.config.HeaderNames.APIKey, c.config.APIKey)
	req.Header.Set(c.config.HeaderNames.UserID, c.config.UserID)
	req.Header.Set(c.config.HeaderNames.DeviceID, sinkerAPIDeviceID)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// TODO: replace with logger
	// fmt.Println("request Headers:", req.Header)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("httpclient do: %v", err)
	}
	defer resp.Body.Close()

	// TODO: replace with logger
	fmt.Println("response Status:", resp.Status)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("io readall: %v", err)
	}

	// TODO: replace with logger
	fmt.Println("response Body:", string(body))

	return body, nil
}
