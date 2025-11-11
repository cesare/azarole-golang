package handlers

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func SignoutHandler(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()

	c.Status(http.StatusOK)
}
