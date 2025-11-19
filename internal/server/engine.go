package server

import (
	"azarole/internal/core"
	"azarole/internal/handlers"
	"azarole/internal/handlers/api"
	"azarole/internal/middlewares"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func Engine(app *core.App) *gin.Engine {
	engine := gin.Default()

	store := cookie.NewStore(app.Secrets.Session.SessionKey.Bytes())
	engine.Use(sessions.Sessions("azarole-session", store))

	engine.Use(cors.New(cors.Config{
		AllowOrigins: []string{app.Config.Frontend.BaseUrl},
		AllowMethods: []string{"DELETE", "GET", "OPTIONS", "POST"},
		AllowHeaders: []string{
			"Content-Type",
		},
		AllowCredentials: true,
	}))

	engine.GET("/ping", handlers.PingHandler)
	engine.DELETE("/signout", handlers.SignoutHandler)
	engine.GET("/current_user", middlewares.RequireSignin(app), handlers.CurrentUserHandler)

	authGroup := engine.Group("/auth/google")
	handlers.RegisterAuthHandlers(authGroup, app)

	apiKeysGroup := engine.Group("/api_keys")
	apiKeysGroup.Use(middlewares.RequireSignin(app))
	handlers.RegisterApiKeysHandlers(*apiKeysGroup, app)

	workplacesGroup := engine.Group("/workplaces")
	workplacesGroup.Use(middlewares.RequireSignin(app))
	handlers.RegisterWorkplacesHandlers(workplacesGroup, app)

	attendancesGroup := engine.Group("/workplaces/:workplace_id/attendance_records")
	attendancesGroup.Use(middlewares.RequireSignin(app))
	handlers.RegisterAttendanceRecordsHandlers(attendancesGroup, app)

	apiGroup := engine.Group("/api")
	apiGroup.Use(middlewares.RequireApiKey(app))
	api.RegisterAttendanceRecordsHandlers(apiGroup, app)

	return engine
}
