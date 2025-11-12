package auth

import (
	"azarole/internal/core"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type AccessTokenRequest struct {
	app *core.App
}

type AccessTokenResponse struct {
	IdToken string `json:"id_token" binding:"required"`
}

func NewAccessTokenRequest(app *core.App) *AccessTokenRequest {
	return &AccessTokenRequest{
		app: app,
	}
}

func (request *AccessTokenRequest) Execute(code string) (*AccessTokenResponse, error) {
	config := request.app.Config
	secrets := request.app.Secrets

	params := url.Values{}
	params.Set("client_id", secrets.GoogleAuth.ClientId)
	params.Set("client_secret", secrets.GoogleAuth.ClientSecret)
	params.Set("code", code)
	params.Set("grant_type", "authorization_code")
	params.Set("redirect_uri", config.Frontend.AuthRedirectUrl())

	resp, err := http.PostForm("https://oauth2.googleapis.com/token", params)
	if err != nil {
		return nil, fmt.Errorf("access token request failed: %s", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read access token response: %s", err)
	}

	var response AccessTokenResponse
	praseError := json.Unmarshal(body, &response)
	if praseError != nil {
		return nil, fmt.Errorf("failed to parse access token response: %s", praseError)
	}

	return &response, nil
}
