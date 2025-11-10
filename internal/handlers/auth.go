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

	type callbackParams struct {
		Code  string `form:"code"`
		State string `form:"state"`
		Error string `form:"error"`
	}

	group.POST("/callback", func(c *gin.Context) {
		var params callbackParams
		err := c.ShouldBind(&params)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		if params.Error != "" {
			c.Status(http.StatusUnauthorized)
			return
		}

		if params.Code == "" || params.State == "" {
			c.Status(http.StatusBadRequest)
			return
		}

		handleSuccess(c, params.Code, params.State)
	})
}

func handleSuccess(c *gin.Context, code string, state string) {
	c.Status(http.StatusNoContent)
}
