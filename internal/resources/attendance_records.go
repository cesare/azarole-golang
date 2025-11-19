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

func (r *AttendanceRecordResource) List(month time.Time) ([]models.AttendandeRecord, error) {
	statement, err := r.app.Database().Prepare("select id, workplace_id, event, recorded_at from attendance_records where workplace_id = $1 and recorded_at >= $2 and recorded_at < $3 order by recorded_at")
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statemnet for listAttendances: %s", err)
	}
	defer statement.Close()

	rows, err := statement.Query(r.workplace.Id, month.UTC(), month.AddDate(0, 1, 0).UTC())
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

func (r *AttendanceRecordResource) Delete(id models.AttendandeRecordId) error {
	statement, err := r.app.Database().Prepare("delete from attendance_records where id = $1 and workplace_id = $2")
	if err != nil {
		return fmt.Errorf("failed to prepare statement for deleteAttendance: %s", err)
	}
	defer statement.Close()

	_, err = statement.Exec(id, r.workplace.Id)
	if err != nil {
		return fmt.Errorf("failed to delete attandance: %s", err)
	}

	return nil
}
