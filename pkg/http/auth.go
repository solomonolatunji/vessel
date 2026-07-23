package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"codedock.dev/codedock/internal/models"
)

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

func (c *Client) Login(email, password string) (*AuthResponse, error) {
	payload := AuthRequest{
		Email:    email,
		Password: password,
	}

	resp, err := c.sendRequest("POST", "/auth/signin", payload)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("login failed: %s", string(respBody))
	}

	var result struct {
		Data *AuthResponse `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

func (c *Client) Logout() error {
	resp, err := c.sendRequest("POST", "/auth/logout", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("logout failed: %s", string(respBody))
	}

	return nil
}
