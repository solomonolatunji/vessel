package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	nethttp "net/http"
	"time"
)

// Client is the Codedock API client.
type Client struct {
	BaseURL    string
	Token      string
	HTTPClient *nethttp.Client
}

// NewClient creates a new API client.
func NewClient(baseURL, token string) *Client {
	return &Client{
		BaseURL: baseURL,
		Token:   token,
		HTTPClient: &nethttp.Client{
			Timeout: time.Second * 30,
		},
	}
}

// sendRequest is a helper for making HTTP requests to the Codedock API.
func (c *Client) sendRequest(method, endpoint string, payload interface{}) (*nethttp.Response, error) {
	var reqBytes []byte
	var err error

	if payload != nil {
		reqBytes, err = json.Marshal(payload)
		if err != nil {
			return nil, err
		}
	}

	req, err := nethttp.NewRequest(method, fmt.Sprintf("%s/api/v1%s", c.BaseURL, endpoint), bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, err
	}

	if c.Token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	}
	req.Header.Set("Content-Type", "application/json")

	return c.HTTPClient.Do(req)
}

// Ping checks if the Codedock server is reachable and the token is valid.
func (c *Client) Ping() error {
	resp, err := c.sendRequest("GET", "/system/health", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusOK {
		return fmt.Errorf("server returned status: %d", resp.StatusCode)
	}

	return nil
}
