package models

type Event struct {
	Base
	Title      string          `json:"title" gorm:"type:varchar(100);not null"`
	Status     string          `json:"status" gorm:"type:varchar(20);not null"`
	Seats      *int            `json:"seats"`
	UserID     uint            `json:"user_id" gorm:"not null;index"`
	CategoryID *uint           `json:"category_id" gorm:"index"`
	Category   *Category       `json:"category" gorm:"foreignKey:CategoryID"`
	Schedule   []EventSchedule `json:"schedule" gorm:"foreignKey:EventID;constraint:OnDelete:CASCADE"`
}
