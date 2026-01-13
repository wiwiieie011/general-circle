package dto

import "time"

type CreateScheduleRequest struct {
	EventID      uint      `json:"event_id"`
	ActivityName string    `json:"activity_name"`
	Speaker      string    `json:"speaker"`
	StartAt      time.Time `json:"start_at"`
	EndAt        time.Time `json:"end_at"`
}

type UpdateScheduleRequest struct {
	EventID      *uint      `json:"event_id"`
	ActivityName *string    `json:"activity_name"`
	Speaker      *string    `json:"speaker"`
	StartAt      *time.Time `json:"start_at"`
	EndAt        *time.Time `json:"end_at"`
}
