package api

import (
	"azarole/internal/core"
	"azarole/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterAttendanceRecordsHandlers(group *gin.RouterGroup, app *core.App) {
	group.POST("/workplaces/:workplace_id/clock_ins", func(c *gin.Context) {
		create(c, models.ClockIn)
	})
	group.POST("/workplaces/:workplace_id/clock_outs", func(c *gin.Context) {
		create(c, models.ClockOut)
	})
}

func create(c *gin.Context, event models.AttendanceEvent) {
	c.Status(http.StatusCreated)
}
