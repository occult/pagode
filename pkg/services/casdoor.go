package services

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/occult/pagode/config"
)

// CasdoorClient wraps the Casdoor SDK client.
type CasdoorClient struct {
	client   *casdoorsdk.Client
	endpoint string
	appName  string
}

// NewCasdoorClient creates a new Casdoor client from configuration.
func NewCasdoorClient(cfg *config.CasdoorConfig) *CasdoorClient {
	client := casdoorsdk.NewClient(
		cfg.Endpoint,
		cfg.ClientId,
		cfg.ClientSecret,
		cfg.Certificate,
		cfg.OrganizationName,
		cfg.ApplicationName,
	)

	return &CasdoorClient{
		client:   client,
		endpoint: cfg.Endpoint,
		appName:  cfg.ApplicationName,
	}
}

// GetSigninURL returns the Casdoor login page URL.
func (c *CasdoorClient) GetSigninURL(redirectURI, state string) string {
	return fmt.Sprintf(
		"%s/login/oauth/authorize?client_id=%s&response_type=code&redirect_uri=%s&scope=openid+profile+email&state=%s",
		c.endpoint,
		url.QueryEscape(c.client.ClientId),
		url.QueryEscape(redirectURI),
		url.QueryEscape(state),
	)
}

// GetSignupURL returns the Casdoor signup page URL.
func (c *CasdoorClient) GetSignupURL(redirectURI, state string) string {
	return fmt.Sprintf(
		"%s/signup/%s?redirect_uri=%s&state=%s",
		c.endpoint,
		url.PathEscape(c.appName),
		url.QueryEscape(redirectURI),
		url.QueryEscape(state),
	)
}

// GetLogoutURL returns the Casdoor logout URL.
func (c *CasdoorClient) GetLogoutURL() string {
	return fmt.Sprintf("%s/api/logout", c.endpoint)
}

// IsReachable checks whether the Casdoor endpoint is reachable.
func (c *CasdoorClient) IsReachable() bool {
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(c.endpoint)
	if err != nil {
		return false
	}
	resp.Body.Close()
	return true
}

// ExchangeCodeAndGetUser exchanges an authorization code for a token and returns the user's email and name.
func (c *CasdoorClient) ExchangeCodeAndGetUser(code, state string) (email, name string, err error) {
	token, err := c.client.GetOAuthToken(code, state)
	if err != nil {
		return "", "", fmt.Errorf("failed to exchange code for token: %w", err)
	}

	claims, err := c.client.ParseJwtToken(token.AccessToken)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse JWT token: %w", err)
	}

	email = claims.Email
	name = claims.Name
	if name == "" {
		name = claims.DisplayName
	}
	if name == "" {
		name = claims.User.Name
	}
	if email == "" {
		return "", "", fmt.Errorf("casdoor token missing email claim")
	}

	return email, name, nil
}
