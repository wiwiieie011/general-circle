package models

import (
	"time"
)

type EventSchedule struct {
	Base
	EventID      uint      `json:"event_id" gorm:"not null;index"`
	Event        Event     `json:"-" gorm:"foreignKey:EventID"`
	ActivityName string    `json:"activity_name" gorm:"type:varchar(100);not null"`
	Speaker      string    `json:"speaker" gorm:"type:varchar(50);not null"`
	StartAt      time.Time `json:"start_at" gorm:"not null"`
	EndAt        time.Time `json:"end_at" gorm:"not null"`
}
