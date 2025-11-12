package handlers

import (
	"azarole/internal/core"
	"azarole/internal/handlers/auth"
	"log/slog"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func RegisterAuthHandlers(group *gin.RouterGroup, application *core.App) {
	group.POST("", func(c *gin.Context) {
		generator := auth.NewAuthorizationRequestGenerator(application)
		authRequest := generator.Generate()

		session := sessions.Default(c)
		session.Set("google-auth-state", authRequest.State)
		session.Set("google-auth-nonce", authRequest.Nonce)
		session.Save()

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

		handleSuccess(c, application, params.Code, params.State)
	})
}

func handleSuccess(c *gin.Context, application *core.App, code string, state string) {
	session := sessions.Default(c)
	savedState, _ := session.Get("google-auth-state").(string)
	savedNonce, _ := session.Get("google-auth-nonce").(string)

	session.Delete("google-auth-state")
	session.Delete("google-auth-nonce")
	session.Save()

	if savedState != state {
		c.Status(http.StatusUnauthorized)
		return
	}

	accessTokenRequest := auth.NewAccessTokenRequest(application)
	accessTokenResponse, err := accessTokenRequest.Execute(code)
	if err != nil {
		slog.Debug("accessTokenRequest failed", "error", err)
		c.Status(http.StatusUnauthorized)
		return
	}

	verifier := auth.NewIdTokenVerifier(application, accessTokenResponse.IdToken, savedNonce)
	claims, err := verifier.Verify()
	if err != nil {
		slog.Debug("verifier failed", "error", err)
		c.Status(http.StatusUnauthorized)
		return
	}

	finder := auth.NewUserFinder(application, claims.Subject)
	result, err := finder.Execute(c)
	if err != nil {
		slog.Debug("finder failed", "error", err)
		c.Status((http.StatusUnauthorized))
		return
	}

	session.Set("userId", result.UserId)
	session.Save()

	c.JSON(http.StatusOK, gin.H{
		"user_id": result.UserId,
	})
}
