package http

import (
	"encoding/json"
	"fmt"
	"io"
	nethttp "net/http"

	"codedock.run/codedock/internal/models"
)

func (c *Client) ListBackups(databaseID string) ([]models.BackupConfig, error) {
	endpoint := "/backups"
	if databaseID != "" {
		endpoint = fmt.Sprintf("/backups?databaseId=%s", databaseID)
	}
	resp, err := c.sendRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to fetch backups (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data []models.BackupConfig `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

func (c *Client) CreateBackup(req interface{}) (*models.BackupConfig, error) {
	resp, err := c.sendRequest("POST", "/backups", req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusCreated && resp.StatusCode != nethttp.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to create backup (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data *models.BackupConfig `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

func (c *Client) TriggerBackup(id string) (*models.BackupRecord, error) {
	resp, err := c.sendRequest("POST", fmt.Sprintf("/backups/%s/trigger", id), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to trigger backup (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data *models.BackupRecord `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

func (c *Client) ListBackupRecords(id string) ([]models.BackupRecord, error) {
	resp, err := c.sendRequest("GET", fmt.Sprintf("/backups/%s/records", id), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to fetch backup records (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data []models.BackupRecord `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, nil
}
