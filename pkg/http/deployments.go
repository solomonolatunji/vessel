package http

import (
	"encoding/json"
	"fmt"
	"io"
	nethttp "net/http"

	"codedock.run/codedock/internal/models"
)

// TriggerDeployment triggers a new manual deployment for an app service.
func (c *Client) TriggerDeployment(serviceID string) (*models.Deployment, error) {
	resp, err := c.sendRequest("POST", fmt.Sprintf("/services/%s/deploy", serviceID), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to trigger deployment (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data *models.Deployment `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

// GetDeploymentStatus checks the status of a specific deployment.
func (c *Client) GetDeploymentStatus(deploymentID string) (*models.Deployment, error) {
	resp, err := c.sendRequest("GET", fmt.Sprintf("/deployments/%s", deploymentID), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to fetch deployment (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data *models.Deployment `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

// ListDeployments gets all deployments for a service.
func (c *Client) ListDeployments(serviceID string) ([]models.Deployment, error) {
	resp, err := c.sendRequest("GET", fmt.Sprintf("/services/%s/deployments", serviceID), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to list deployments (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data []models.Deployment `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

// GetDeploymentLogs fetches build/run logs for a specific deployment.
func (c *Client) GetDeploymentLogs(deploymentID string) (string, error) {
	resp, err := c.sendRequest("GET", fmt.Sprintf("/deployments/%s/logs", deploymentID), nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to fetch logs (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data string `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Data, nil
}

// GetServiceMetrics fetches performance metrics for a service.
func (c *Client) GetServiceMetrics(serviceID string) ([]models.ServiceMetric, error) {
	resp, err := c.sendRequest("GET", fmt.Sprintf("/services/%s/metrics", serviceID), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to fetch metrics (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data []models.ServiceMetric `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, nil
}
