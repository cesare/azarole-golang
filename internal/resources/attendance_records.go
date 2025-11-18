package resources

import (
	"azarole/internal/core"
	"azarole/internal/models"
	"fmt"
	"time"
)

type AttendanceRecordResource struct {
	app       *core.App
	workplace *models.Workplace
}

func NewAttendanceRecordResource(app *core.App, workplace *models.Workplace) *AttendanceRecordResource {
	return &AttendanceRecordResource{
		app:       app,
		workplace: workplace,
	}
}

func (r *AttendanceRecordResource) Create(event models.AttendanceEvent) (*models.AttendandeRecord, error) {
	statement, err := r.app.Database().Prepare("insert into attendance_records (workplace_id, event, recorded_at, created_at) values ($1, $2, $3, $4) returning id, workplace_id, event, recorded_at")
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement for findWorkplaces: %s", err)
	}
	defer statement.Close()

	now := time.Now().UTC()
	rows, err := statement.Query(r.workplace.Id, event, now, now)
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
