package handlers

import (
	"azarole/internal/core"
	"azarole/internal/models"
	"azarole/internal/views"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterApiKeysHandlers(group gin.RouterGroup, app *core.App) {
	group.GET("", func(c *gin.Context) {
		currentUser := c.MustGet("currentUser").(models.User)
		apiKeys, err := findApiKeys(app, &currentUser)
		if err != nil {
			slog.Debug("failed to find apiKeys", "error", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		ks := []views.ApiKeyView{}
		for _, apiKey := range apiKeys {
			view := views.FromApiKey(&apiKey)
			ks = append(ks, *view)
		}

		c.JSON(http.StatusOK, gin.H{
			"api_keys": ks, // TODO: fix key into camelCase
		})
	})
}

func findApiKeys(app *core.App, user *models.User) ([]models.ApiKey, error) {
	statement, err := app.Database().Prepare("select id, user_id, name, digest, created_at from api_keys where user_id = $1 order by created_at desc")
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement for findApiKeys: %s", err)
	}
	defer statement.Close()

	rows, err := statement.Query(user.Id)
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
