package auth

import (
	"azarole/internal/core"
	"fmt"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
)

type IdTokenClaims struct {
	Nonce string `json:"nonce"`
	jwt.RegisteredClaims
}

type IdTokenVerifier struct {
	application *core.App
	token       string
	nonce       string
}

func NewIdTokenVerifier(application *core.App, token string, nonce string) *IdTokenVerifier {
	return &IdTokenVerifier{
		application: application,
		token:       token,
		nonce:       nonce,
	}
}

func (verifier *IdTokenVerifier) Verify() (*IdTokenClaims, error) {
	k, err := keyfunc.NewDefault([]string{"https://www.googleapis.com/oauth2/v3/certs"})
	if err != nil {
		return nil, fmt.Errorf("failed to build keyfunc: %s", err)
	}

	parsed, err := jwt.ParseWithClaims(verifier.token, &IdTokenClaims{}, k.Keyfunc)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %s", err)
	}

	claims := parsed.Claims.(*IdTokenClaims)
	if claims.Nonce != verifier.nonce {
		return nil, fmt.Errorf("nonces don't match")
	}

	return claims, nil
}
