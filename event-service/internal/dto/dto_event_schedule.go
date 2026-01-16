package dto

import "time"

type CreateScheduleRequest struct {
	ActivityName string    `json:"activity_name" binding:"required,min=3,max=100"`
	Speaker      string    `json:"speaker" binding:"required,min=3,max=50"`
	StartAt      time.Time `json:"start_at" binding:"required"`
	EndAt        time.Time `json:"end_at" binding:"required"`
}

type UpdateScheduleRequest struct {
	ActivityName *string    `json:"activity_name"`
	Speaker      *string    `json:"speaker"`
	StartAt      *time.Time `json:"start_at"`
	EndAt        *time.Time `json:"end_at"`
}
