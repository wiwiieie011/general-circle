package models

import (
	"time"

	"gorm.io/gorm"
)

type EventSchedule struct {
	gorm.Model
	EventID      uint      `json:"event_id" gorm:"uniqueIndex"`
	ActivityName string    `json:"activity_name"`
	Speaker      string    `json:"speaker"`
	StartAt      time.Time `json:"start_at"`
	EndAt        time.Time `json:"end_at"`
}
