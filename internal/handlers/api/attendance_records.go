package api

import (
	"azarole/internal/core"
	"azarole/internal/models"
	"azarole/internal/resources"
	"azarole/internal/views"
	"fmt"
	"net/http"
	"time"

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

	attendance, err := createAttendance(app, workplace, event)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	view := views.FromAttendanceRecord(attendance)
	c.JSON(http.StatusCreated, gin.H{
		"attendanceRecord": view,
	})
}

func createAttendance(app *core.App, workplace *models.Workplace, event models.AttendanceEvent) (*models.AttendandeRecord, error) {
	statement, err := app.Database().Prepare("insert into attendance_records (workplace_id, event, recorded_at, created_at) values ($1, $2, $3, $4) returning id, workplace_id, event, recorded_at")
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement for findWorkplaces: %s", err)
	}
	defer statement.Close()

	now := time.Now().UTC()
	rows, err := statement.Query(workplace.Id, event, now, now)
	if err != nil {
		return nil, fmt.Errorf("query for creating attandance record failed: %s", err)
	}
	defer rows.Close()

	rows.Next()

	var ar models.AttendandeRecord
	err = rows.Scan(&ar.Id, &ar.WorkplaceId, &ar.Event, &ar.RecordedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to map row into attendanceRecord: %s", err)
	}

	return &ar, nil
}
