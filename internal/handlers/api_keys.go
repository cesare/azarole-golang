package handlers

import (
	"azarole/internal/core"
	apikeys "azarole/internal/handlers/api_keys"
	"azarole/internal/models"
	"azarole/internal/resources"
	"azarole/internal/views"
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

		rs := resources.NewApiKeyResources(app, &currentUser)
		apiKeys, err := rs.List()
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
			"apiKeys": ks,
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
			"apiKey": details,
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

		rs := resources.NewApiKeyResources(app, &currentUser)
		rs.Delete(apiKeyId)

		c.Status(http.StatusOK)
	})
}
