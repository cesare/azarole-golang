package handlers

import (
	app "azarole/internal"
	"azarole/internal/handlers/auth"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func RegisterAuthHandlers(group *gin.RouterGroup, application *app.Application) {
	group.POST("", func(c *gin.Context) {
		generator := auth.NewAuthorizationRequestGenerator(application)
		authRequest := generator.Generate()

		session := sessions.Default(c)
		session.Set("google-auth-state", authRequest.State)
		session.Set("google-auth-nonce", authRequest.Nonce)

		c.JSON(http.StatusOK, gin.H{
			"location": authRequest.RequestUrl,
		})
	})

	group.POST("/callback", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})
}
