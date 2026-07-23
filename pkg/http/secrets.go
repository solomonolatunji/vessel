package http

import (
	"encoding/json"
	"fmt"
	"io"
	nethttp "net/http"

	"codedock.dev/codedock/internal/models"
)

// GetSecrets retrieves the secrets/environment variables for a project.
func (c *Client) GetSecrets(projectID string) (models.VarsRequest, error) {
	resp, err := c.sendRequest("GET", fmt.Sprintf("/projects/%s/env", projectID), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get secrets (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data models.VarsRequest `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

// SetSecrets updates the secrets/environment variables for a project.
func (c *Client) SetSecrets(projectID string, req models.SetEnvVarsRequest) error {
	resp, err := c.sendRequest("PUT", fmt.Sprintf("/projects/%s/env", projectID), req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to set secrets (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}
