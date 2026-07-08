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
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/solomonolatunji/vessel/internal/store"
	"github.com/solomonolatunji/vessel/internal/types"
)

// GitService handles Git provider API integration, repository listing, and automated local cloning and pulling.
type GitService struct {
	store      *store.Store
	httpClient *http.Client
}

// NewGitService initializes a GitService wired to the underlying store and an HTTP client.
func NewGitService(s *store.Store) *GitService {
	return &GitService{
		store: s,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// SaveProvider stores a user's GitHub or GitLab access token and account name.
func (gs *GitService) SaveProvider(userID string, req *types.GitConnectRequest) (*types.GitProviderConfig, error) {
	switch req.Provider {
	case "github", "gitlab":
	default:
		return nil, errors.New("unsupported git provider; must be 'github' or 'gitlab'")
	}
	if req.AccessToken == "" {
		return nil, errors.New("access token is required")
	}

	gp := &types.GitProviderConfig{
		UserID:      userID,
		Provider:    req.Provider,
		AccessToken: req.AccessToken,
		AccountName: req.AccountName,
	}
	if err := gs.store.SaveGitProvider(gp); err != nil {
		return nil, fmt.Errorf("failed to save git provider: %w", err)
	}
	gp.AccessToken = ""
	return gp, nil
}

// GetConnectedProviders returns the connection status for GitHub and GitLab for a specific user.
func (gs *GitService) GetConnectedProviders(userID string) ([]map[string]any, error) {
	var results []map[string]any

	for _, provider := range []string{"github", "gitlab"} {
		gp, err := gs.store.GetGitProvider(userID, provider)
		if err != nil {
			return nil, err
		}
		if gp != nil {
			results = append(results, map[string]any{
				"provider":    provider,
				"connected":   true,
				"accountName": gp.AccountName,
				"updatedAt":   gp.UpdatedAt,
			})
		} else {
			results = append(results, map[string]any{
				"provider":  provider,
				"connected": false,
			})
		}
	}
	return results, nil
}

// DisconnectProvider removes a stored OAuth/PAT connection for a specific user and provider.
func (gs *GitService) DisconnectProvider(userID, provider string) error {
	return gs.store.DeleteGitProvider(userID, provider)
}

// ListRepositories retrieves public and private repositories for the user from GitHub or GitLab API.
func (gs *GitService) ListRepositories(ctx context.Context, userID, provider string) ([]types.GitRepository, error) {
	gp, err := gs.store.GetGitProvider(userID, provider)
	if err != nil {
		return nil, fmt.Errorf("failed to load git credentials: %w", err)
	}
	if gp == nil || gp.AccessToken == "" {
		return nil, fmt.Errorf("user is not authenticated with %s", provider)
	}

	switch provider {
	case "github":
		return gs.listGitHubRepos(ctx, gp.AccessToken)
	case "gitlab":
		return gs.listGitLabRepos(ctx, gp.AccessToken)
	default:
		return nil, errors.New("unsupported provider: " + provider)
	}
}

func (gs *GitService) listGitHubRepos(ctx context.Context, token string) ([]types.GitRepository, error) {
	reqURL := "https://api.github.com/user/repos?per_page=100&sort=updated"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := gs.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("github api request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("github api returned status %d: %s", resp.StatusCode, string(body))
	}

	var ghRepos []struct {
		ID       int64  `json:"id"`
		Name     string `json:"name"`
		FullName string `json:"full_name"`
		Private  bool   `json:"private"`
		CloneURL string `json:"clone_url"`
		HTMLURL  string `json:"html_url"`
		Default  string `json:"default_branch"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&ghRepos); err != nil {
		return nil, fmt.Errorf("failed to decode github repositories: %w", err)
	}

	var results []types.GitRepository
	for _, r := range ghRepos {
		results = append(results, types.GitRepository{
			ID:            r.ID,
			Name:          r.Name,
			FullName:      r.FullName,
			Private:       r.Private,
			CloneURL:      r.CloneURL,
			HTMLURL:       r.HTMLURL,
			DefaultBranch: r.Default,
		})
	}
	return results, nil
}

func (gs *GitService) listGitLabRepos(ctx context.Context, token string) ([]types.GitRepository, error) {
	reqURL := "https://gitlab.com/api/v4/projects?membership=true&per_page=100&order_by=updated_at"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := gs.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("gitlab api request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("gitlab api returned status %d: %s", resp.StatusCode, string(body))
	}

	var glRepos []struct {
		ID         int64  `json:"id"`
		Name       string `json:"name"`
		FullName   string `json:"path_with_namespace"`
		Visibility string `json:"visibility"`
		CloneURL   string `json:"http_url_to_repo"`
		HTMLURL    string `json:"web_url"`
		Default    string `json:"default_branch"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&glRepos); err != nil {
		return nil, fmt.Errorf("failed to decode gitlab projects: %w", err)
	}

	var results []types.GitRepository
	for _, r := range glRepos {
		results = append(results, types.GitRepository{
			ID:            r.ID,
			Name:          r.Name,
			FullName:      r.FullName,
			Private:       r.Visibility == "private",
			CloneURL:      r.CloneURL,
			HTMLURL:       r.HTMLURL,
			DefaultBranch: r.Default,
		})
	}
	return results, nil
}

// CloneOrPullRepository checks out or updates the codebase from the project's RepositoryURL into targetDir.
func (gs *GitService) CloneOrPullRepository(ctx context.Context, project *types.ProjectConfig, targetDir string, logWriter io.Writer) error {
	repoURL := strings.TrimSpace(project.RepositoryURL)
	if repoURL == "" {
		return errors.New("repositoryUrl is not set for project")
	}

	branch := strings.TrimSpace(project.Branch)
	if branch == "" {
		branch = "main"
	}

	authURL := gs.injectAuthTokenIfAvailable(repoURL)

	if logWriter != nil {
		fmt.Fprintf(logWriter, "📥 [GitService] Preparing to sync codebase from %s (branch: %s)...\n", repoURL, branch)
	}

	gitDir := filepath.Join(targetDir, ".git")
	if _, err := os.Stat(gitDir); err == nil {
		if logWriter != nil {
			fmt.Fprintf(logWriter, "🔄 [GitService] Existing local directory detected; pulling latest changes...\n")
		}
		fetchCmd := exec.CommandContext(ctx, "git", "-C", targetDir, "fetch", "origin", branch)
		if out, err := fetchCmd.CombinedOutput(); err != nil {
			return fmt.Errorf("git fetch failed: %v (%s)", err, string(out))
		}
		resetCmd := exec.CommandContext(ctx, "git", "-C", targetDir, "reset", "--hard", "origin/"+branch)
		if out, err := resetCmd.CombinedOutput(); err != nil {
			return fmt.Errorf("git reset failed: %v (%s)", err, string(out))
		}
		if logWriter != nil {
			fmt.Fprintf(logWriter, "✅ [GitService] Successfully updated local repository to latest commit on %s.\n", branch)
		}
		return nil
	}

	_ = os.RemoveAll(targetDir)
	if err := os.MkdirAll(filepath.Dir(targetDir), 0755); err != nil {
		return fmt.Errorf("failed to create build parent dir: %w", err)
	}

	var cloneCmd *exec.Cmd
	cloneArgs := []string{"clone", "--depth", "1", "-b", branch, authURL, targetDir}
	if logWriter != nil {
		fmt.Fprintf(logWriter, "🚀 [GitService] Running git clone --depth 1 -b %s...\n", branch)
	}

	cloneCmd = exec.CommandContext(ctx, "git", cloneArgs...)
	var stderr bytes.Buffer
	cloneCmd.Stderr = &stderr
	if err := cloneCmd.Run(); err != nil {
		if strings.Contains(stderr.String(), "Remote branch") && branch == "main" {
			if logWriter != nil {
				fmt.Fprintf(logWriter, "⚠️ [GitService] Branch 'main' not found; retrying clone with repository default branch...\n")
			}
			_ = os.RemoveAll(targetDir)
			cloneCmd = exec.CommandContext(ctx, "git", "clone", "--depth", "1", authURL, targetDir)
			if errFallback := cloneCmd.Run(); errFallback != nil {
				return fmt.Errorf("git clone failed: %v (%s)", errFallback, stderr.String())
			}
		} else {
			return fmt.Errorf("git clone failed: %v (%s)", err, stderr.String())
		}
	}

	if logWriter != nil {
		fmt.Fprintf(logWriter, "✅ [GitService] Successfully cloned repository into %s.\n", targetDir)
	}
	return nil
}

func (gs *GitService) injectAuthTokenIfAvailable(repoURL string) string {
	u, err := url.Parse(repoURL)
	if err != nil || u.Scheme != "https" {
		return repoURL
	}

	var provider string
	if strings.Contains(u.Host, "github.com") {
		provider = "github"
	} else if strings.Contains(u.Host, "gitlab.com") {
		provider = "gitlab"
	} else {
		return repoURL
	}

	gp, err := gs.store.GetGitProvider("", provider)
	if err != nil || gp == nil || gp.AccessToken == "" {
		return repoURL
	}

	switch provider {
	case "github":
		u.User = url.UserPassword("x-access-token", gp.AccessToken)
	case "gitlab":
		u.User = url.UserPassword("oauth2", gp.AccessToken)
	}
	return u.String()
}
