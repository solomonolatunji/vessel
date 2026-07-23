package http

import (
	"encoding/json"
	"fmt"
	"io"
	nethttp "net/http"

	"codedock.dev/codedock/internal/models"
)

// ListEnvironments returns all environments for a given project ID.
func (c *Client) ListEnvironments(projectID string) ([]*models.EnvironmentConfig, error) {
	resp, err := c.sendRequest("GET", fmt.Sprintf("/projects/%s/environments", projectID), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to list environments (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data []*models.EnvironmentConfig `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

// CreateEnvironment creates a new environment for a project.
func (c *Client) CreateEnvironment(projectID string, req *models.EnvironmentConfig) (*models.EnvironmentConfig, error) {
	resp, err := c.sendRequest("POST", fmt.Sprintf("/projects/%s/environments", projectID), req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusCreated && resp.StatusCode != nethttp.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to create environment (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data *models.EnvironmentConfig `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

// DeleteEnvironment deletes an environment by its ID.
func (c *Client) DeleteEnvironment(id string) error {
	resp, err := c.sendRequest("DELETE", fmt.Sprintf("/environments/%s", id), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusOK && resp.StatusCode != nethttp.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete environment (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}
