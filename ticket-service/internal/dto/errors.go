package dto

import "errors"

var (
	ErrEventNotFound     = errors.New("event not found")
	ErrEventNotPublished = errors.New("event not published")
	ErrEventNotStarted   = errors.New("event not started")
	ErrEventEnded        = errors.New("event already ended")
	ErrTicketSoldOut     = errors.New("event already ended")
)
