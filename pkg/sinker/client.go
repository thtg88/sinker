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
	config 			*config.SinkerAPI
	httpClient	*http.Client
	logger			*log.Logger
}

func NewAPIClient(httpClient *http.Client, cfg *config.SinkerAPI, logger *log.Logger) *APIClient {
	return &APIClient{
		config:			cfg,
		httpClient:	httpClient,
		logger:			logger,
	}
}

func (c *APIClient) sinkerApiRequest(method string, uri string, requestBody []byte, sinkerAPIDeviceID string) ([]byte, error) {
	url := fmt.Sprint(c.config.BaseURL, uri)

	c.logger.Printf("sinkerapirequest url: %s", url)

	req, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("http newrequest: %v", err)
	}

	req.Header.Set(c.config.HeaderNames.APIKey, c.config.APIKey)
	req.Header.Set(c.config.HeaderNames.UserID, c.config.UserID)
	req.Header.Set(c.config.HeaderNames.DeviceID, sinkerAPIDeviceID)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// c.logger.Printf("sinkerapirequest request headers: %v", req.Header)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("httpclient do: %v", err)
	}
	defer resp.Body.Close()

	c.logger.Printf("sinkerapirequest response status: %d", resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("io readall: %v", err)
	}

	c.logger.Printf("sinkerapirequest response body: %s", string(body))

	return body, nil
}
