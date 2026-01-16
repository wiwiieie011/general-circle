package dto


type UpdateNotificationPreferencesRequest struct {
	TicketPurchased *bool `json:"ticket_purchased"`
	EventCanceled   *bool `json:"event_canceled"`
	EventReminder   *bool `json:"event_reminder"`
	PushEnabled     *bool `json:"push_enabled"`
	InAppEnabled    *bool `json:"in_app_enabled"`
}

type NotificationType string

const (
	NotificationTypeTicket   NotificationType = "ticket_purchase"
	NotificationTypeEvent    NotificationType = "event_notification"
	NotificationTypeReminder NotificationType = "reminder"
)

type TicketPurchasedEvent struct {
	UserID     uint   `json:"user_id"`
	EventID    uint   `json:"event_id"`
	EventTitle string `json:"event_title"`
}

type EventCancelledEvent struct {
	EventID    uint   `json:"event_id"`
	EventTitle string `json:"event_title"`
	UserIDs    []uint `json:"user_ids"` // всех владельцев билетов
}

type EventReminder struct {
	EventID    uint   `json:"event_id"`
	EventTitle string `json:"event_title"`
	UserIDs    []uint `json:"user_ids"` // всех владельцев билетов
}


