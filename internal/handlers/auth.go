package handlers

import (
	app "azarole/internal"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterAuthHandlers(group *gin.RouterGroup, application *app.Application) {
	group.POST("", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	group.POST("/callback", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})
}
