package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"

	"vessel.dev/vessel/internal/models"
	"vessel.dev/vessel/internal/repositories"
)

type GitAppsService struct {
	repo repositories.GitAppRepository
}

func NewGitAppsService(repo repositories.GitAppRepository) *GitAppsService {
	return &GitAppsService{repo: repo}
}

func listApps[T any](ctx context.Context, teamID string, listFn func(context.Context, string) ([]T, error)) ([]T, error) {
	if teamID == "" {
		return nil, errors.New("team ID is required")
	}
	return listFn(ctx, teamID)
}

func getApp[T any](ctx context.Context, id string, getFn func(context.Context, string) (*T, error)) (*T, error) {
	if id == "" {
		return nil, errors.New("app ID is required")
	}
	return getFn(ctx, id)
}

func saveApp[T any](ctx context.Context, app *T, getTeamID func(*T) string, getID func(*T) string, setID func(*T, string), saveFn func(context.Context, *T) error) error {
	if app == nil {
		return errors.New("app config is required")
	}
	if getTeamID(app) == "" {
		return errors.New("team ID is required")
	}
	if getID(app) == "" {
		setID(app, uuid.NewString())
	}
	return saveFn(ctx, app)
}

func deleteApp(ctx context.Context, id string, deleteFn func(context.Context, string) error) error {
	if id == "" {
		return errors.New("app ID is required")
	}
	return deleteFn(ctx, id)
}

type githubManifestConversionResponse struct {
	ID            int    `json:"id"`
	Slug          string `json:"slug"`
	ClientID      string `json:"client_id"`
	ClientSecret  string `json:"client_secret"`
	WebhookSecret string `json:"webhook_secret"`
	PEM           string `json:"pem"`
	HTMLURL       string `json:"html_url"`
	Name          string `json:"name"`
}

func (s *GitAppsService) ExchangeGithubManifestCode(ctx context.Context, code string, teamID string) (*models.GithubApp, error) {
	if code == "" {
		return nil, errors.New("conversion code is required")
	}
	if teamID == "" {
		teamID = "default"
	}

	url := fmt.Sprintf("https://api.github.com/app-manifests/%s/conversions", code)
	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("github api error: %s", string(body))
	}

	var conversion githubManifestConversionResponse
	if err := json.NewDecoder(resp.Body).Decode(&conversion); err != nil {
		return nil, err
	}

	app := &models.GithubApp{
		ID:            uuid.NewString(),
		TeamID:        teamID,
		Name:          conversion.Name,
		AppID:         fmt.Sprintf("%d", conversion.ID),
		ClientID:      conversion.ClientID,
		ClientSecret:  conversion.ClientSecret,
		WebhookSecret: conversion.WebhookSecret,
		PrivateKey:    conversion.PEM,
		IsPublic:      false,
	}

	if err := s.repo.SaveGithubApp(ctx, app); err != nil {
		return nil, err
	}

	return app, nil
}

func (s *GitAppsService) ListGithubApps(ctx context.Context, teamID string) ([]models.GithubApp, error) {
	return listApps(ctx, teamID, s.repo.ListGithubApps)
}

func (s *GitAppsService) GetGithubApp(ctx context.Context, id string) (*models.GithubApp, error) {
	return getApp(ctx, id, s.repo.GetGithubApp)
}

func (s *GitAppsService) SaveGithubApp(ctx context.Context, app *models.GithubApp) error {
	return saveApp(ctx, app, func(a *models.GithubApp) string { return a.TeamID }, func(a *models.GithubApp) string { return a.ID }, func(a *models.GithubApp, id string) { a.ID = id }, s.repo.SaveGithubApp)
}

func (s *GitAppsService) DeleteGithubApp(ctx context.Context, id string) error {
	return deleteApp(ctx, id, s.repo.DeleteGithubApp)
}

func (s *GitAppsService) ListGitlabApps(ctx context.Context, teamID string) ([]models.GitlabApp, error) {
	return listApps(ctx, teamID, s.repo.ListGitlabApps)
}

func (s *GitAppsService) GetGitlabApp(ctx context.Context, id string) (*models.GitlabApp, error) {
	return getApp(ctx, id, s.repo.GetGitlabApp)
}

func (s *GitAppsService) SaveGitlabApp(ctx context.Context, app *models.GitlabApp) error {
	return saveApp(ctx, app, func(a *models.GitlabApp) string { return a.TeamID }, func(a *models.GitlabApp) string { return a.ID }, func(a *models.GitlabApp, id string) { a.ID = id }, s.repo.SaveGitlabApp)
}

func (s *GitAppsService) DeleteGitlabApp(ctx context.Context, id string) error {
	return deleteApp(ctx, id, s.repo.DeleteGitlabApp)
}

func (s *GitAppsService) ListBitbucketApps(ctx context.Context, teamID string) ([]models.BitbucketApp, error) {
	return listApps(ctx, teamID, s.repo.ListBitbucketApps)
}

func (s *GitAppsService) GetBitbucketApp(ctx context.Context, id string) (*models.BitbucketApp, error) {
	return getApp(ctx, id, s.repo.GetBitbucketApp)
}

func (s *GitAppsService) SaveBitbucketApp(ctx context.Context, app *models.BitbucketApp) error {
	return saveApp(ctx, app, func(a *models.BitbucketApp) string { return a.TeamID }, func(a *models.BitbucketApp) string { return a.ID }, func(a *models.BitbucketApp, id string) { a.ID = id }, s.repo.SaveBitbucketApp)
}

func (s *GitAppsService) DeleteBitbucketApp(ctx context.Context, id string) error {
	return deleteApp(ctx, id, s.repo.DeleteBitbucketApp)
}
