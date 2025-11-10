package auth

import (
	app "azarole/internal"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/url"
)

type AuthorizationRequest struct {
	State      string
	Nonce      string
	RequestUrl string
}

type AuthorizationRequestGenerator struct {
	application *app.Application
}

func NewAuthorizationRequestGenerator(application *app.Application) *AuthorizationRequestGenerator {
	return &AuthorizationRequestGenerator{
		application: application,
	}
}

func (generator *AuthorizationRequestGenerator) Generate() *AuthorizationRequest {
	state := generator.generateRandomString()
	nonce := generator.generateRandomString()
	requestUrl := generator.generateRequestUrl(state, nonce)

	return &AuthorizationRequest{
		State:      state,
		Nonce:      nonce,
		RequestUrl: requestUrl,
	}
}

func (generator *AuthorizationRequestGenerator) generateRandomString() string {
	bytes := make([]byte, 36)
	rand.Read(bytes)

	return base64.RawStdEncoding.EncodeToString(bytes)
}

func (generator *AuthorizationRequestGenerator) generateRequestUrl(state string, nonce string) string {
	clientId := generator.application.Secrets.GoogleAuth.ClientId
	callbackUrl := generator.callbackUrl()

	params := url.Values{}
	params.Set("client_id", clientId)
	params.Set("redirect_uri", callbackUrl)
	params.Set("response_type", "code")
	params.Set("scope", "openid email")
	params.Set("state", state)
	params.Set("nonce", nonce)

	requestUrl := url.URL{
		Scheme:   "https",
		Host:     "accounts.google.com",
		Path:     "/o/oauth2/v2/auth",
		RawQuery: params.Encode(),
	}
	return requestUrl.String()
}

func (generator *AuthorizationRequestGenerator) callbackUrl() string {
	baseUrl := generator.application.Config.Frontend.BaseUrl
	return fmt.Sprintf("%s/signin/callback", baseUrl)
}
