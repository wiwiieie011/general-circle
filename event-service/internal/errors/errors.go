package errors

import "errors"

var (
	ErrEventIsNil             = errors.New("event is nil")
	ErrEventScheduleIsNil     = errors.New("event schedule is nil")
	ErrCategoryIsNil          = errors.New("category is nil")
	ErrEmptyTitle             = errors.New("title cannot be empty")
	ErrCategoryNotFound       = errors.New("category not found")
	ErrEventNotFound          = errors.New("event not found")
	ErrEventScheduleNotFound  = errors.New("event schedule not found")
	ErrEventIsNotDraft        = errors.New("you can delete or publish only draft status event")
	ErrEventIsNotPublished    = errors.New("you can cancel only published status event")
	ErrEmptyName              = errors.New("name cannot be empty")
	ErrCategoryNameExists     = errors.New("category already has this name")
	ErrEmptyActivityName      = errors.New("activity name cannot be empty")
	ErrNotCorrectID           = errors.New("id cannot be less than 1")
	ErrEmptySpeaker           = errors.New("speaker cannot be empty")
	ErrNotCorrectScheduleTime = errors.New("start time cannot be equal and after end time and vice versa")
	ErrNotCorrectNum          = errors.New("number cannot be less than 1")
)
