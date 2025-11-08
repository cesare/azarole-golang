package app

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func Engine(config *Config, secrets *Secrets) (*gin.Engine, error) {
	engine := gin.Default()

	store := cookie.NewStore(secrets.Session.SessionKey.Bytes())
	engine.Use(sessions.Sessions("azarole-session", store))

	engine.Use(cors.New(cors.Config{
		AllowOrigins: []string{config.Frontend.BaseUrl},
		AllowMethods: []string{"DELETE", "GET", "OPTIONS", "POST"},
		AllowHeaders: []string{
			"Content-Type",
		},
		AllowCredentials: true,
	}))

	return engine, nil
}
