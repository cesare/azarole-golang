package resources

import (
	"azarole/internal/core"
	"azarole/internal/models"
	"fmt"
)

type ApiKeyResources struct {
	app  *core.App
	user *models.User
}

func NewApiKeyResources(app *core.App, user *models.User) *ApiKeyResources {
	return &ApiKeyResources{
		app:  app,
		user: user,
	}
}

func (r *ApiKeyResources) List() ([]models.ApiKey, error) {
	statement, err := r.app.Database().Prepare("select id, user_id, name, digest, created_at from api_keys where user_id = $1 order by created_at desc")
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement for findApiKeys: %s", err)
	}
	defer statement.Close()

	rows, err := statement.Query(r.user.Id)
	if err != nil {
		return nil, fmt.Errorf("query for findApiKeys failed: %s", err)
	}
	defer rows.Close()

	apiKeys := []models.ApiKey{}
	for rows.Next() {
		var apiKey models.ApiKey
		err = rows.Scan(&apiKey.Id, &apiKey.UserId, &apiKey.Name, &apiKey.Digest, &apiKey.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to map row into apiKey: %s", err)
		}

		apiKeys = append(apiKeys, apiKey)
	}

	return apiKeys, nil
}

func (r *ApiKeyResources) Delete(apiKeyId models.ApiKeyId) error {
	statement, err := r.app.Database().Prepare("delete from api_keys where user_id = $1 and id = $2")
	if err != nil {
		return fmt.Errorf("failed to prepare statement for deleteApiKey: %s", err)
	}
	defer statement.Close()

	_, err = statement.Exec(r.user.Id, apiKeyId)
	if err != nil {
		return fmt.Errorf("failed to delete api key: %s", err)
	}

	return nil
}
