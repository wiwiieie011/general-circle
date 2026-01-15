package dto

import "errors"

var (
	ErrEventNotFound     = errors.New("event not found")
	ErrEventNotPublished = errors.New("event not published")
)
