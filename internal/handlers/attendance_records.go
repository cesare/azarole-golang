package handlers

import (
	"azarole/internal/core"
	"azarole/internal/models"
	"azarole/internal/resources"
	"azarole/internal/views"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func RegisterAttendanceRecordsHandlers(group *gin.RouterGroup, app *core.App) {
	group.GET("", func(c *gin.Context) {
		currentUser := c.MustGet("currentUser").(models.User)

		v := c.Param("workplace_id")
		workplaceId, err := models.FromStringToWorkplaceId(v)
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}

		params := listingParams{}
		err = c.ShouldBind(&params)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		params = *params.Normalize()

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

		targetMonth := params.ToTime()
		attendances, err := listAttendances(app, workplace, targetMonth)
		if err != nil {
			slog.Debug("Failed to listing attendances", "error", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		vs := []views.AttendandeRecordView{}
		for _, a := range attendances {
			v := views.FromAttendanceRecord(&a)
			vs = append(vs, *v)
		}

		c.JSON(http.StatusOK, gin.H{
			"year":              params.Year,
			"month":             params.Month,
			"workplace":         views.FromWorkplace(workplace),
			"attendanceRecords": vs,
		})
	})
}

type listingParams struct {
	Year  int `form:"year"`
	Month int `form:"month"`
}

func (p *listingParams) Normalize() *listingParams {
	now := time.Now()
	year := p.Year
	if year == 0 {
		year = now.Year()
	}

	month := p.Month
	if month == 0 {
		month = int(now.Month())
	}

	quotient := month / 12
	if quotient > 0 {
		month = month % 12
		year += quotient
	}

	return &listingParams{
		Year:  year,
		Month: month,
	}
}

func (p *listingParams) ToTime() time.Time {
	now := time.Now()
	location := now.Location()
	return time.Date(p.Year, time.Month(p.Month), 1, 0, 0, 0, 0, location)
}

func listAttendances(app *core.App, workplace *models.Workplace, month time.Time) ([]models.AttendandeRecord, error) {
	statement, err := app.Database().Prepare("select id, workplace_id, event, recorded_at from attendance_records where workplace_id = $1 and recorded_at >= $2 and recorded_at < $3 order by recorded_at")
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statemnet for listAttendances: %s", err)
	}
	defer statement.Close()

	rows, err := statement.Query(workplace.Id, month.UTC(), month.AddDate(0, 1, 0).UTC())
	if err != nil {
		return nil, fmt.Errorf("query failed: %s", err)
	}
	defer rows.Close()

	as := []models.AttendandeRecord{}
	for rows.Next() {
		var a models.AttendandeRecord
		err = rows.Scan(&a.Id, &a.WorkplaceId, &a.Event, &a.RecordedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to map row into AttandanceRecord: %s", err)
		}
		as = append(as, a)
	}

	return as, nil
}
