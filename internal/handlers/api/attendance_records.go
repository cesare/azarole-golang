package api

import (
	"azarole/internal/core"
	"azarole/internal/models"
	"azarole/internal/resources"
	"azarole/internal/views"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterAttendanceRecordsHandlers(group *gin.RouterGroup, app *core.App) {
	group.POST("/workplaces/:workplace_id/clock_ins", func(c *gin.Context) {
		create(c, app, models.ClockIn)
	})
	group.POST("/workplaces/:workplace_id/clock_outs", func(c *gin.Context) {
		create(c, app, models.ClockOut)
	})
}

func create(c *gin.Context, app *core.App, event models.AttendanceEvent) {
	currentUser := c.MustGet("currentUser").(models.User)

	v := c.Param("workplace_id")
	workplaceId, err := models.FromStringToWorkplaceId(v)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	wrs := resources.NewWorkplaceResources(app, &currentUser)
	workplace, err := wrs.Find(workplaceId)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	if workplace == nil {
		c.Status(http.StatusNotFound)
		return
	}

	ars := resources.NewAttendanceRecordResource(app, workplace)
	attendance, err := ars.Create(event)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	view := views.FromAttendanceRecord(attendance)
	c.JSON(http.StatusCreated, gin.H{
		"attendanceRecord": view,
	})
}
