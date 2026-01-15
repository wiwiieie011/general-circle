package models

import "gorm.io/gorm"

type Event struct {
	gorm.Model
	Title      string         `json:"title"`
	Status     string         `json:"status"`
	Seats      *int           `json:"seats"`
	UserID     uint           `json:"user_id"`
	CategoryID *uint          `json:"category_id"`
	Category   *Category      `json:"category" gorm:"foreignKey:CategoryID"`
	Schedule   *EventSchedule `json:"schedule" gorm:"foreignKey:EventID;constraint:OnDelete:CASCADE"`
}
