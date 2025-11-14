package core

import (
	"encoding/base64"
	"fmt"

	"github.com/caarlos0/env/v11"
)

type b64string struct {
	bytes []byte
}

func (b64 *b64string) Bytes() []byte {
	return b64.bytes
}

func (b64 *b64string) UnmarshalText(text []byte) error {
	bytes := make([]byte, len(text))
	_, err := base64.StdEncoding.Decode(bytes, text)
	if err != nil {
		return fmt.Errorf("failed to decode base64 value: %s", err)
	}

	b64.bytes = bytes
	return nil
}

type ApiKeyConfig struct {
	DigestingSecretKey b64string `env:"API_KEY_DIGESTING_SECRET_KEY,required"`
}

func loadApiKey() (*ApiKeyConfig, error) {
	cfg, err := env.ParseAs[ApiKeyConfig]()
	if err != nil {
		return nil, fmt.Errorf("failed to load ApiKeyConfig: %s", err)
	}

	return &cfg, nil
}

type GoogleAuthConfig struct {
	ClientId     string `env:"GOOGLE_AUTH_CLIENT_ID,required"`
	ClientSecret string `env:"GOOGLE_AUTH_CLIENT_SECRET,required"`
}

func loadGoogleAuth() (*GoogleAuthConfig, error) {
	cfg, err := env.ParseAs[GoogleAuthConfig]()
	if err != nil {
		return nil, fmt.Errorf("failed to load GoogleAuthConfig: %s", err)
	}

	return &cfg, nil
}

type SessionConfig struct {
	SessionKey b64string `env:"SESSION_KEY,required"`
}

func loadSession() (*SessionConfig, error) {
	cfg, err := env.ParseAs[SessionConfig]()
	if err != nil {
		return nil, fmt.Errorf("failed to load SessionConfig: %s", err)
	}

	return &cfg, nil
}

type Secrets struct {
	ApiKey     ApiKeyConfig
	GoogleAuth GoogleAuthConfig
	Session    SessionConfig
}

func LoadSecrets() (*Secrets, error) {
	apiKey, err := loadApiKey()
	if err != nil {
		return nil, err
	}

	googleAuth, err := loadGoogleAuth()
	if err != nil {
		return nil, err
	}

	session, err := loadSession()
	if err != nil {
		return nil, err
	}

	secrets := Secrets{
		ApiKey:     *apiKey,
		GoogleAuth: *googleAuth,
		Session:    *session,
	}
	return &secrets, nil
}
