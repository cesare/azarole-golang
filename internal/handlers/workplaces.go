package handlers

import (
	"azarole/internal/core"
	"azarole/internal/models"
	"azarole/internal/resources"
	"azarole/internal/views"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createWorkplaceParams struct {
	Name string `form:"name" binding:"required"`
}

func RegisterWorkplacesHandlers(group *gin.RouterGroup, app *core.App) {
	group.GET("", func(c *gin.Context) {
		currentUser := c.MustGet("currentUser").(models.User)

		rs := resources.NewWorkplaceResources(app, &currentUser)
		workplaces, err := rs.List()
		if err != nil {
			slog.Debug("failed to find workplaces", "error", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		vs := []views.WorkplaceView{}
		for _, wp := range workplaces {
			vs = append(vs, *views.FromWorkplace(&wp))
		}

		c.JSON(http.StatusOK, gin.H{
			"workplaces": vs,
		})
	})

	group.POST("", func(c *gin.Context) {
		currentUser := c.MustGet("currentUser").(models.User)

		var params createWorkplaceParams
		err := c.ShouldBind(&params)
		if err != nil {
			slog.Debug("failed to bind parameters", "error", err)
			c.Status(http.StatusBadRequest)
			return
		}

		rs := resources.NewWorkplaceResources(app, &currentUser)
		workplace, err := rs.Create(params.Name)
		if err != nil {
			slog.Debug("failed to create workplace", "error", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		view := views.FromWorkplace(workplace)
		c.JSON(http.StatusCreated, gin.H{
			"workplace": view,
		})
	})
}
