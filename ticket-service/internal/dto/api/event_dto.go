package dto_api

type EventStatus string

const (
	EventStatusPublished  EventStatus = "published"
)

type EventResponse struct {
	Status EventStatus `json:"status"`
}
