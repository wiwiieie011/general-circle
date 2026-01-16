package dto

import "time"

type CreateScheduleRequest struct {
	ActivityName string    `json:"activity_name"`
	Speaker      string    `json:"speaker"`
	StartAt      time.Time `json:"start_at"`
	EndAt        time.Time `json:"end_at"`
}

type UpdateScheduleRequest struct {
	ActivityName *string    `json:"activity_name"`
	Speaker      *string    `json:"speaker"`
	StartAt      *time.Time `json:"start_at"`
	EndAt        *time.Time `json:"end_at"`
}
