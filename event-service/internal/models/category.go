package models

type Category struct {
	Base
	Name   string  `json:"name" gorm:"type:varchar(50);not null;uniqueIndex"`
	Events []Event `json:"-" gorm:"foreignKey:CategoryID"`
}
