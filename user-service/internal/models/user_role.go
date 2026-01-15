package models

import "time"

type UserRole string

const (
	RoleUser      UserRole = "user"
	RoleOrganizer UserRole = "organizer"
	RoleAdmin     UserRole = "admin"
)

type User struct {
	ID        uint      `gorm:"primaryKey"`
	Email     string    `gorm:"uniqueIndex;not null"`
	Password  string    `gorm:"not null"`
	FirstName string
	LastName  string
	Role      UserRole  `gorm:"type:varchar(20);not null;default:'user'"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
