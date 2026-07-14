package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Client is the Vessl API client.
type Client struct {
	BaseURL    string
	Token      string
	HTTPClient *http.Client
}

// NewClient creates a new API client.
func NewClient(baseURL, token string) *Client {
	return &Client{
		BaseURL: baseURL,
		Token:   token,
		HTTPClient: &http.Client{
			Timeout: time.Second * 30,
		},
	}
}

// sendRequest is a helper for making HTTP requests to the Vessl API.
func (c *Client) sendRequest(method, endpoint string, payload interface{}) (*http.Response, error) {
	var reqBytes []byte
	var err error

	if payload != nil {
		reqBytes, err = json.Marshal(payload)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, fmt.Sprintf("%s/api/v1%s", c.BaseURL, endpoint), bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, err
	}

	if c.Token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	}
	req.Header.Set("Content-Type", "application/json")

	return c.HTTPClient.Do(req)
}

// Ping checks if the Vessl server is reachable and the token is valid.
func (c *Client) Ping() error {
	resp, err := c.sendRequest("GET", "/system/health", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned status: %d", resp.StatusCode)
	}

	return nil
}
