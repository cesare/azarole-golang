package apikeys

import (
	"azarole/internal/core"
	"azarole/internal/models"
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
	return nil, nil
}
