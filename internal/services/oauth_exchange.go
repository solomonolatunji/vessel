package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"vessel.dev/vessel/internal/models"
)

func ExchangeCode(p *models.OAuthProviderConfig, code string) (string, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	switch strings.ToLower(p.ProviderName) {
	case "github":
		return exchangeGitHub(client, p, code)
	case "bitbucket":
		return exchangeBitbucket(client, p, code)
	default:
		return exchangeOIDC(client, p, code)
	}
}

func GetAuthorizationURL(p *models.OAuthProviderConfig, state string) (string, error) {
	if !p.Enabled || p.ClientID == "" {
		return "", fmt.Errorf("oauth provider %s is not enabled or configured", p.ProviderName)
	}
	switch strings.ToLower(p.ProviderName) {
	case "github":
		return fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&scope=user:email&state=%s",
			url.QueryEscape(p.ClientID), url.QueryEscape(p.RedirectURI), url.QueryEscape(state)), nil
	case "gitlab":
		base := p.BaseURL
		if base == "" {
			base = "https://gitlab.com"
		}
		return fmt.Sprintf("%s/oauth/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=read_user+openid+profile+email&state=%s",
			strings.TrimRight(base, "/"), url.QueryEscape(p.ClientID), url.QueryEscape(p.RedirectURI), url.QueryEscape(state)), nil
	case "google":
		return fmt.Sprintf("https://accounts.google.com/o/oauth2/v2/auth?client_id=%s&redirect_uri=%s&response_type=code&scope=openid+email+profile&state=%s",
			url.QueryEscape(p.ClientID), url.QueryEscape(p.RedirectURI), url.QueryEscape(state)), nil
	case "azuread":
		tenant := p.Tenant
		if tenant == "" {
			tenant = "common"
		}
		return fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=openid+email+profile&state=%s",
			url.PathEscape(tenant), url.QueryEscape(p.ClientID), url.QueryEscape(p.RedirectURI), url.QueryEscape(state)), nil
	case "discord":
		return fmt.Sprintf("https://discord.com/api/oauth2/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=identify+email&state=%s",
			url.QueryEscape(p.ClientID), url.QueryEscape(p.RedirectURI), url.QueryEscape(state)), nil
	case "authentik", "zitadel", "clerk", "infomaniak":
		if p.BaseURL == "" {
			return "", fmt.Errorf("base url is required for %s oauth", p.ProviderName)
		}
		authEndpoint := strings.TrimRight(p.BaseURL, "/") + "/oauth/authorize"
		switch strings.ToLower(p.ProviderName) {
		case "authentik":
			authEndpoint = strings.TrimRight(p.BaseURL, "/") + "/application/o/authorize/"
		case "zitadel":
			authEndpoint = strings.TrimRight(p.BaseURL, "/") + "/oauth/v2/authorize"
		}
		return fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&response_type=code&scope=openid+email+profile&state=%s",
			authEndpoint, url.QueryEscape(p.ClientID), url.QueryEscape(p.RedirectURI), url.QueryEscape(state)), nil
	case "bitbucket":
		return fmt.Sprintf("https://bitbucket.org/site/oauth2/authorize?client_id=%s&response_type=code&state=%s",
			url.QueryEscape(p.ClientID), url.QueryEscape(state)), nil
	default:
		return "", fmt.Errorf("unsupported oauth provider: %s", p.ProviderName)
	}
}

func executeOAuthFlow(client *http.Client, req *http.Request, userURL string, parseEmail func(io.Reader) (string, error)) (string, error) {
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var tokenResp struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil || tokenResp.AccessToken == "" {
		return "", fmt.Errorf("failed to get access token")
	}
	userReq, _ := http.NewRequest("GET", userURL, nil)
	userReq.Header.Set("Authorization", "Bearer "+tokenResp.AccessToken)
	userResp, err := client.Do(userReq)
	if err != nil {
		return "", err
	}
	defer userResp.Body.Close()
	return parseEmail(userResp.Body)
}

func exchangeGitHub(client *http.Client, p *models.OAuthProviderConfig, code string) (string, error) {
	body, _ := json.Marshal(map[string]string{
		"client_id": p.ClientID, "client_secret": p.ClientSecret,
		"code": code, "redirect_uri": p.RedirectURI,
	})
	req, _ := http.NewRequest("POST", "https://github.com/login/oauth/access_token", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	return executeOAuthFlow(client, req, "https://api.github.com/user/emails", func(r io.Reader) (string, error) {
		var emails []struct {
			Email   string `json:"email"`
			Primary bool   `json:"primary"`
		}
		if err := json.NewDecoder(r).Decode(&emails); err == nil {
			for _, e := range emails {
				if e.Primary {
					return e.Email, nil
				}
			}
			if len(emails) > 0 {
				return emails[0].Email, nil
			}
		}
		return "", fmt.Errorf("could not retrieve email from github")
	})
}

func exchangeBitbucket(client *http.Client, p *models.OAuthProviderConfig, code string) (string, error) {
	values := url.Values{
		"client_id": {p.ClientID}, "client_secret": {p.ClientSecret},
		"code": {code}, "grant_type": {"authorization_code"},
	}
	req, _ := http.NewRequest("POST", "https://bitbucket.org/site/oauth2/access_token", strings.NewReader(values.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return executeOAuthFlow(client, req, "https://api.bitbucket.org/2.0/user/emails", func(r io.Reader) (string, error) {
		var emailResp struct {
			Values []struct {
				Email     string `json:"email"`
				IsPrimary bool   `json:"is_primary"`
			} `json:"values"`
		}
		if err := json.NewDecoder(r).Decode(&emailResp); err == nil {
			for _, e := range emailResp.Values {
				if e.IsPrimary {
					return e.Email, nil
				}
			}
		}
		return "", fmt.Errorf("could not retrieve email from bitbucket")
	})
}

func exchangeOIDC(client *http.Client, p *models.OAuthProviderConfig, code string) (string, error) {
	tokenURL, userURL := oidcEndpoints(p)
	values := url.Values{
		"client_id": {p.ClientID}, "client_secret": {p.ClientSecret},
		"code": {code}, "grant_type": {"authorization_code"}, "redirect_uri": {p.RedirectURI},
	}
	req, _ := http.NewRequest("POST", tokenURL, strings.NewReader(values.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return executeOAuthFlow(client, req, userURL, func(r io.Reader) (string, error) {
		var userInfo struct {
			Email string `json:"email"`
		}
		b, _ := io.ReadAll(r)
		if err := json.Unmarshal(b, &userInfo); err == nil && userInfo.Email != "" {
			return userInfo.Email, nil
		}
		return "", fmt.Errorf("could not extract email from %s user info", p.ProviderName)
	})
}

func oidcEndpoints(p *models.OAuthProviderConfig) (tokenURL, userURL string) {
	switch strings.ToLower(p.ProviderName) {
	case "gitlab":
		base := p.BaseURL
		if base == "" {
			base = "https://gitlab.com"
		}
		return strings.TrimRight(base, "/") + "/oauth/token",
			strings.TrimRight(base, "/") + "/api/v4/user"
	case "discord":
		return "https://discord.com/api/oauth2/token", "https://discord.com/api/users/@me"
	case "azuread":
		tenant := p.Tenant
		if tenant == "" {
			tenant = "common"
		}
		return fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", tenant),
			"https://graph.microsoft.com/oidc/userinfo"
	case "authentik":
		return strings.TrimRight(p.BaseURL, "/") + "/application/o/token/",
			strings.TrimRight(p.BaseURL, "/") + "/application/o/userinfo/"
	case "zitadel":
		return strings.TrimRight(p.BaseURL, "/") + "/oauth/v2/token",
			strings.TrimRight(p.BaseURL, "/") + "/oidc/v1/userinfo"
	default:
		return "https://oauth2.googleapis.com/token",
			"https://openidconnect.googleapis.com/v1/userinfo"
	}
}
