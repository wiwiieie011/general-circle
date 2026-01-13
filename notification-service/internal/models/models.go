package models

import (
	"time"

	"gorm.io/gorm"
)

type Notification struct {
	gorm.Model
	UserID    uint `json:"user_id" gorm:"index"`
	Type      string `json:"type"` // тут либо покупка билетов, уведомление о мероприятиях, и напоминания 
	Title     string `json:"title"`
	Body      string `json:"body"`
	Read      bool `json:"read"`
	CreatedAt time.Time `json:"created_at"`
}

type NotificationPreference struct {
	UserID uint `gorm:"primaryKey"`

	TicketPurchased bool // отключает уведомление о покупке билетов 
	EventCanceled   bool // отклюает уведомления о мероприятиях
	EventReminder   bool // отключает напоминания 

	PushEnabled  bool   
	InAppEnabled bool
}