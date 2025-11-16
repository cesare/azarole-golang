package views

import (
	"azarole/internal/models"
	"time"
)

type AttendandeRecordView struct {
	Id          models.AttendandeRecordId `json:"id"`
	WorkplaceId models.WorkplaceId        `json:"workplace_id"`
	Event       models.AttendanceEvent    `json:"event"`
	RecordedAt  time.Time                 `json:"recordedAt"`
}

func FromAttendanceRecord(a *models.AttendandeRecord) *AttendandeRecordView {
	return &AttendandeRecordView{
		Id:         a.Id,
		Event:      a.Event,
		RecordedAt: a.RecordedAt,
	}
}
