package models

import (
	"time"

	"gorm.io/gorm"
)

type Notification struct {
	Model
	UserID  uint   `json:"user_id" gorm:"index"`
	EventID uint   `json:"event_id" gorm:"index"`
	Type    string `json:"type"` // тут либо покупка билетов, уведомление о мероприятиях, и напоминания
	Title   string `json:"title"`
	Body    string `json:"body"`
	Read    bool   `json:"read"`
}

type NotificationPreference struct {
	UserID uint `gorm:"primaryKey"`

	TicketPurchased bool // отключает уведомление о покупке билетов
	EventCanceled   bool // отключает уведомления о мероприятиях
	EventReminder   bool // отключает напоминания

	PushEnabled  bool
	InAppEnabled bool
}

type Model struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
