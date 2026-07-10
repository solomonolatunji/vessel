package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"vessel.dev/vessel/internal/models"
	"vessel.dev/vessel/internal/repositories"
)

type VercelService struct {
	vercelRepo *repositories.VercelRepository
}

func NewVercelService(repo *repositories.VercelRepository) *VercelService {
	return &VercelService{vercelRepo: repo}
}

func (s *VercelService) GetAccount(ctx context.Context, userID string, teamID *string) (*models.UserVercelAccount, error) {
	return s.vercelRepo.GetAccount(ctx, userID, teamID)
}

func (s *VercelService) GetAccountsForUser(ctx context.Context, userID string) ([]*models.UserVercelAccount, error) {
	return s.vercelRepo.GetAccountsForUser(ctx, userID)
}

func (s *VercelService) HandleCallback(ctx context.Context, userID, code string) (*models.UserVercelAccount, error) {
	clientID := os.Getenv("VERCEL_CLIENT_ID")
	clientSecret := os.Getenv("VERCEL_CLIENT_SECRET")
	redirectURI := os.Getenv("VERCEL_REDIRECT_URI")

	if clientID == "" || clientSecret == "" {
		return nil, errors.New("vercel oauth is not configured on this server")
	}

	data := url.Values{}
	data.Set("code", code)
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("redirect_uri", redirectURI)

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.vercel.com/v2/oauth/access_token", bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to exchange code for vercel token: %s", string(body))
	}

	var result struct {
		AccessToken string `json:"access_token"`
		TeamID      string `json:"team_id"`
		UserID      string `json:"user_id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	account := &models.UserVercelAccount{
		UserID:      userID,
		AccessToken: result.AccessToken,
		AccountName: "Vercel Account", // We could fetch user info to get the real name, but this is a default
	}

	if result.TeamID != "" {
		account.TeamID = &result.TeamID
	}

	if err := s.vercelRepo.SaveAccount(ctx, account); err != nil {
		return nil, err
	}
	return account, nil
}

func (s *VercelService) ListProjects(ctx context.Context, userID string, teamID *string) ([]models.VercelProject, error) {
	account, err := s.GetAccount(ctx, userID, teamID)
	if err != nil || account == nil {
		return nil, errors.New("vercel account not linked")
	}

	reqURL := "https://api.vercel.com/v9/projects"
	if account.TeamID != nil {
		reqURL += "?teamId=" + *account.TeamID
	}

	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+account.AccessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to list vercel projects: %s", string(body))
	}

	var result struct {
		Projects []struct {
			ID          string `json:"id"`
			Name        string `json:"name"`
			Framework   string `json:"framework"`
			NodeVersion string `json:"nodeVersion"`
			AccountID   string `json:"accountId"`
		} `json:"projects"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	var projects []models.VercelProject
	for _, p := range result.Projects {
		projects = append(projects, models.VercelProject{
			ID:          p.ID,
			Name:        p.Name,
			Framework:   p.Framework,
			NodeVersion: p.NodeVersion,
			AccountID:   p.AccountID,
		})
	}
	return projects, nil
}

func (s *VercelService) GetProjectEnvVars(ctx context.Context, userID string, teamID *string, projectID string) ([]models.VercelEnvVar, error) {
	account, err := s.GetAccount(ctx, userID, teamID)
	if err != nil || account == nil {
		return nil, errors.New("vercel account not linked")
	}

	reqURL := fmt.Sprintf("https://api.vercel.com/v9/projects/%s/env", projectID)
	if account.TeamID != nil {
		reqURL += "?teamId=" + *account.TeamID
	}

	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+account.AccessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to list vercel env vars: %s", string(body))
	}

	var result struct {
		Envs []models.VercelEnvVar `json:"envs"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Envs, nil
}
