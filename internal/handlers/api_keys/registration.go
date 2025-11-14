package apikeys

import (
	"azarole/internal/core"
	"azarole/internal/models"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"
)

type RegistrationDetails struct {
	Id    models.ApiKeyId `json:"id"`
	Name  string          `json:"name"`
	Token string          `json:"token"`
}

type ApiKeyRegistration struct {
	app  *core.App
	user *models.User
	name string
}

func NewApiKeyRegistration(app *core.App, user *models.User, name string) *ApiKeyRegistration {
	return &ApiKeyRegistration{
		app:  app,
		user: user,
		name: name,
	}
}

func (r *ApiKeyRegistration) Execute() (*RegistrationDetails, error) {
	token := r.generateToken()
	digest := r.digestToken(token)
	apiKey, err := r.createApiKey(digest)
	if err != nil {
		return nil, err
	}

	details := RegistrationDetails{
		Id:    apiKey.Id,
		Name:  apiKey.Name,
		Token: token,
	}
	return &details, nil
}

func (r *ApiKeyRegistration) generateToken() string {
	bytes := make([]byte, 96)
	rand.Read(bytes)

	return base64.RawURLEncoding.EncodeToString(bytes)
}

func (r *ApiKeyRegistration) digestToken(token string) string {
	secretKey := r.app.Secrets.ApiKey.DigestingSecretKey.Bytes()
	h := hmac.New(sha256.New, secretKey)
	h.Write([]byte(token))
	digest := h.Sum(nil)

	return hex.EncodeToString(digest)
}

func (r *ApiKeyRegistration) createApiKey(digest string) (*models.ApiKey, error) {
	statement, err := r.app.Database().Prepare("insert into api_keys (user_id, name, digest, created_at) values ($1, $2, $3, $4) returning id, user_id, name, digest, created_at")
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement for createApiKey: %s", err)
	}
	defer statement.Close()

	var apiKey models.ApiKey
	now := time.Now().UTC()
	err = statement.QueryRow(r.user.Id, r.name, digest, now).Scan(&apiKey.Id, &apiKey.UserId, &apiKey.Name, &apiKey.Digest, &apiKey.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("query for createApiKey failed: %s", err)
	}

	return &apiKey, nil
}
