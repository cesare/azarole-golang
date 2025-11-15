package models

import (
	"time"
)

type AttendandeRecordId uint32

type AttendanceEvent string

const (
	ClockIn  AttendanceEvent = "clock-in"
	ClockOut AttendanceEvent = "clock-out"
)

type AttendandeRecord struct {
	Id          AttendandeRecordId
	WorkplaceId WorkplaceId
	Event       AttendanceEvent
	RecordedAt  time.Time
}
