package middlewares

import (
	"azarole/internal/core"
	"azarole/internal/models"
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func RequireSignin(app *core.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userId, ok := session.Get("userId").(models.UserId)
		if !ok {
			c.Status(http.StatusUnauthorized)
			return
		}

		user, err := loadCurrentUser(app, userId)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		if user == nil {
			c.Status(http.StatusUnauthorized)
			return
		}

		c.Set("currentUser", user)
		c.Next()
	}
}

func loadCurrentUser(app *core.App, userId models.UserId) (*models.User, error) {
	statement, err := app.Database().Prepare("select id from users where id = $1")
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement for loadCurrentUser: %s", err)
	}
	defer statement.Close()

	rows, err := statement.Query(userId)
	if err != nil {
		return nil, fmt.Errorf("query for loadCurrentUser failed: %s", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	var user models.User
	rows.Scan(&user.Id)

	return &user, nil
}
