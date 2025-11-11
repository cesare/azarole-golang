package auth

import (
	app "azarole/internal"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type AccessTokenRequest struct {
	appplication *app.Application
}

type AccessTokenResponse struct {
	IdToken string `json:"id_token" binding:"required"`
}

func NewAccessTokenRequest(application *app.Application) *AccessTokenRequest {
	return &AccessTokenRequest{
		appplication: application,
	}
}

func (request *AccessTokenRequest) Execute(code string) (*AccessTokenResponse, error) {
	config := request.appplication.Config
	secrets := request.appplication.Secrets

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
