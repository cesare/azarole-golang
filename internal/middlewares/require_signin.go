package middlewares

import (
	"azarole/internal/core"
	"azarole/internal/models"
	"azarole/internal/resources"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func RequireSignin(app *core.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userId, ok := session.Get("userId").(uint32)
		if !ok {
			c.Status(http.StatusUnauthorized)
			c.Abort()
			return
		}

		rs := resources.NewUserResources(app)
		user, err := rs.Find(models.UserId(userId))
		if err != nil {
			c.Status(http.StatusInternalServerError)
			c.Abort()
			return
		}
		if user == nil {
			c.Status(http.StatusUnauthorized)
			c.Abort()
			return
		}

		c.Set("currentUser", *user)
		c.Next()
	}
}
