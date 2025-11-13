package handlers

import (
	"azarole/internal/models"
	"azarole/internal/views"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CurrentUserHandler(c *gin.Context) {
	user := c.MustGet("currentUser").(models.User)
	view := views.FromUser(&user)

	c.JSON(http.StatusOK, gin.H{
		"user": view,
	})
}
