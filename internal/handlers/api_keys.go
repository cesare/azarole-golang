package handlers

import (
	"azarole/internal/core"
	apikeys "azarole/internal/handlers/api_keys"
	"azarole/internal/models"
	"azarole/internal/views"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createApiKeyParams struct {
	Name string `form:"name" binding:"required"`
}

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

	group.POST("", func(c *gin.Context) {
		currentUser := c.MustGet("currentUser").(models.User)

		var params createApiKeyParams
		err := c.ShouldBind(&params)
		if err != nil {
			slog.Debug("failed to bind createApiKeyParams", "error", err)
			c.Status(http.StatusBadRequest)
			return
		}

		registration := apikeys.NewApiKeyRegistration(app, &currentUser, params.Name)
		details, err := registration.Execute()
		if err != nil {
			slog.Debug("api key registration failed", "error", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"api_key": details, // TODO fix key into camelCase
		})
	})

	group.DELETE("/:id", func(c *gin.Context) {
		currentUser := c.MustGet("currentUser").(models.User)

		v := c.Param("id")
		apiKeyId, err := models.FromStringToApiKeyId(v)
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}

		deleteApiKey(app, &currentUser, apiKeyId)

		c.Status(http.StatusOK)
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

func deleteApiKey(app *core.App, user *models.User, apiKeyId models.ApiKeyId) error {
	statement, err := app.Database().Prepare("delete from api_keys where user_id = $1 and id = $2")
	if err != nil {
		return fmt.Errorf("failed to prepare statement for deleteApiKey: %s", err)
	}
	defer statement.Close()

	_, err = statement.Exec(user.Id, apiKeyId)
	if err != nil {
		return fmt.Errorf("failed to delete api key: %s", err)
	}

	return nil
}
