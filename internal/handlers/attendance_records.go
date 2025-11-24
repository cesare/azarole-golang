package handlers

import (
	"azarole/internal/core"
	"azarole/internal/models"
	"azarole/internal/resources"
	"azarole/internal/views"
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
		ars := resources.NewAttendanceRecordResource(app, workplace)
		attendances, err := ars.List(targetMonth)
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

	type deletingPath struct {
		WorkplaceId models.WorkplaceId        `uri:"workplace_id"`
		Id          models.AttendandeRecordId `uri:"id"`
	}

	group.DELETE("/:id", func(c *gin.Context) {
		currentUser := c.MustGet("currentUser").(models.User)

		var path deletingPath
		err := c.BindUri(&path)
		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}

		wrs := resources.NewWorkplaceResources(app, &currentUser)
		workplace, err := wrs.Find(path.WorkplaceId)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		if workplace == nil {
			c.Status(http.StatusNotFound)
			return
		}

		ars := resources.NewAttendanceRecordResource(app, workplace)
		err = ars.Delete(path.Id)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Status(http.StatusOK)
	})
}

type listingParams struct {
	Year  int `form:"year"`
	Month int `form:"month"`
}

func (p *listingParams) Normalize() *listingParams {
	defaultLocation := time.FixedZone("Asia/Tokyo", 9*60*60)
	now := time.Now().In(defaultLocation)
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
