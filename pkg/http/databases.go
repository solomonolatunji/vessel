package http

import (
	"encoding/json"
	"fmt"
	"io"
	nethttp "net/http"

	"codedock.dev/codedock/internal/models"
)

// ListDatabases returns all databases.
func (c *Client) ListDatabases(projectID string) ([]*models.Database, error) {
	url := "/databases"
	if projectID != "" {
		url = fmt.Sprintf("/databases?projectId=%s", projectID)
	}
	resp, err := c.sendRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to list databases (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data []*models.Database `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

// GetDatabase retrieves a single database by its ID.
func (c *Client) GetDatabase(id string) (*models.Database, error) {
	resp, err := c.sendRequest("GET", fmt.Sprintf("/databases/%s", id), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get database (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data *models.Database `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

// CreateDatabase creates a new database.
func (c *Client) CreateDatabase(req *models.CreateDatabaseRequest) (*models.Database, error) {
	resp, err := c.sendRequest("POST", "/databases", req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusCreated && resp.StatusCode != nethttp.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to create database (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data *models.Database `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

// DeleteDatabase deletes a database by its ID.
func (c *Client) DeleteDatabase(id string) error {
	resp, err := c.sendRequest("DELETE", fmt.Sprintf("/databases/%s", id), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusOK && resp.StatusCode != nethttp.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete database (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

// ImportDatabase triggers a data import from a source URL into the database.
func (c *Client) ImportDatabase(id string, req *models.ImportDatabaseRequest) error {
	resp, err := c.sendRequest("POST", fmt.Sprintf("/databases/%s/import", id), req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to import database (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}
