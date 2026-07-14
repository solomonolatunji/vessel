package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"vessl.dev/vessl/internal/models"
)

func (c *Client) Me() (*models.User, error) {
	resp, err := c.sendRequest("GET", "/auth/me", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get profile: %s", string(body))
	}

	var result struct {
		Data *models.User `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result.Data, nil
}
