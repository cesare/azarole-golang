package server

import (
	"azarole/internal/core"
	"azarole/internal/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func Engine(application *core.Application) *gin.Engine {
	engine := gin.Default()

	store := cookie.NewStore(application.Secrets.Session.SessionKey.Bytes())
	engine.Use(sessions.Sessions("azarole-session", store))

	engine.Use(cors.New(cors.Config{
		AllowOrigins: []string{application.Config.Frontend.BaseUrl},
		AllowMethods: []string{"DELETE", "GET", "OPTIONS", "POST"},
		AllowHeaders: []string{
			"Content-Type",
		},
		AllowCredentials: true,
	}))

	engine.GET("/ping", handlers.PingHandler)
	engine.DELETE("/signout", handlers.SignoutHandler)

	authGroup := engine.Group("/auth/google")
	handlers.RegisterAuthHandlers(authGroup, application)

	return engine
}
