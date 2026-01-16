package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"notification-service/internal/dto"
	"notification-service/internal/models"
	"notification-service/internal/services"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	brokers []string
	srv     services.NotificationService
	log     *slog.Logger
	groupID string
	topics  []string
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewConsumer(brokers []string, srv services.NotificationService, log *slog.Logger) *Consumer {
	ctx, cancel := context.WithCancel(context.Background())
	return &Consumer{
		brokers: brokers,
		srv:     srv,
		log:     log,
		groupID: "notification-service",
		topics:  []string{"ticket.purchased", "event.cancelled", "event.reminder"},
		ctx:     ctx,
		cancel:  cancel,
	}
}

func (c *Consumer) Start() {
	for _, topic := range c.topics {
		go c.consumeTopic(topic)
	}
}

func (c *Consumer) consumeTopic(topic string) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  c.brokers,
		GroupID:  c.groupID,
		Topic:    topic,
		MinBytes: 1,
		MaxBytes: 10e6,
	})

	defer r.Close()

	for {
		m, err := r.ReadMessage(c.ctx)
		if err != nil {
			c.log.Warn("failed to read message:", err)
			continue
		}

		c.log.Info("Received message from topic %s: %s\n", topic, string(m.Value))

		switch topic {
		case "ticket.purchased":
			c.handleTicketPurchased(m.Value)
		case "event.cancelled":
			c.handleEventCancelled(m.Value)
		case "event.reminder":
			c.handleEventReminder(m.Value)
		}
	}

}

func (c *Consumer) handleTicketPurchased(payload []byte) {
	var evt dto.TicketPurchasedEvent
	if err := json.Unmarshal(payload, &evt); err != nil {
		c.log.Error("failed to unmarshal ticket purchased:", err)
		return
	}

	pref, err := c.srv.GetNotificationPreferences(evt.UserID)
	if err != nil {
		c.log.Error("failed to load preferences", "user_id", evt.UserID, "error", err)
		return
	}

	if !pref.TicketPurchased {
		c.log.Info(
			"ticket purchased notification disabled",
			"user_id", evt.UserID,
		)
		return
	}

	notification := &models.Notification{
		UserID:  evt.UserID,
		EventID: evt.EventID,
		Type:    string(dto.NotificationTypeTicket),
		Title:   "Билет куплен",
		Body:    fmt.Sprintf("Ты успешно приобрёл билет на %s", evt.EventTitle),
	}

	if err := c.srv.CreateNotificationInternal(notification); err != nil {
		c.log.Error("failed to create notification:", err)
	}
}

func (c *Consumer) handleEventCancelled(payload []byte) {
	var evt dto.EventCancelledEvent
	if err := json.Unmarshal(payload, &evt); err != nil {
		c.log.Error("failed to unmarshal event reminder:", err)
		return
	}

	for _, userID := range evt.UserIDs {

		pref, err := c.srv.GetNotificationPreferences(userID)
		if err != nil {
			c.log.Error("failed to load preferences", "user_id", userID)
			continue
		}

		if !pref.EventCanceled {
			continue
		}

		notification := &models.Notification{
			UserID:  userID,
			EventID: evt.EventID,
			Type:    string(dto.NotificationTypeEvent),
			Title:   "Мероприятие отменено",
			Body:    fmt.Sprintf("Мероприятие %s отменено", evt.EventTitle),
		}
		if err := c.srv.CreateNotificationInternal(notification); err != nil {
			c.log.Error("failed to create  notification:", err)
		}
	}
}
func (c *Consumer) handleEventReminder(payload []byte) {
	var evt dto.EventReminder
	if err := json.Unmarshal(payload, &evt); err != nil {
		c.log.Error("failed to unmarshal event cancelled:", err)
		return
	}

	for _, userID := range evt.UserIDs {

		pref, err := c.srv.GetNotificationPreferences(userID)
		if err != nil {
			c.log.Error("failed to load preferences", "user_id", userID)
			continue
		}

		if !pref.EventReminder {
			continue
		}
		notification := &models.Notification{
			UserID:  userID,
			EventID: evt.EventID,
			Type:    string(dto.NotificationTypeReminder),
			Title:   "Напоминание о мероприятии",
			Body:    fmt.Sprintf("Завтра состоится мероприятие %s", evt.EventTitle),
		}
		if err := c.srv.CreateNotificationInternal(notification); err != nil {
			c.log.Error("failed to create  notification:", err)
		}
	}
}

func (c *Consumer) Stop() {
	c.cancel()
}
