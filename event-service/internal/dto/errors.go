package dto

import "errors"

var (
	ErrEventIsNil         = errors.New("event is nil")
	ErrEventScheduleIsNil = errors.New("event schedule is nil")
	ErrCategoryIsNil      = errors.New("category is nil")
)
