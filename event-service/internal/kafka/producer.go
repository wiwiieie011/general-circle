package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	eventCancelled = "event.cancelled"
	eventReminder  = "event.reminder"
)

type Producer struct {
	writer *kafka.Writer
	logger *slog.Logger
}

type EventProducer interface {
	SendEventCancelled(ctx context.Context, eventID uint) error
	SendEventReminder(ctx context.Context, eventID uint, eventTitle string, eventDate time.Time) error
	Close() error
}

type EventCancelledMessage struct {
	EventID uint `json:"event_id"`
}

type EventReminderMessage struct {
	EventID    uint      `json:"event_id"`
	EventTitle string    `json:"event_title"`
	EventDate  time.Time `json:"event_date"`
}

func NewProducer(brokers []string, logger *slog.Logger) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Balancer:     &kafka.LeastBytes{},
			WriteTimeout: 20 * time.Second,
			ReadTimeout:  10 * time.Second,
			RequiredAcks: kafka.RequireOne,
			Async:        false,
			Compression:  kafka.Snappy,
		},
		logger: logger,
	}
}

func (p *Producer) SendEventCancelled(ctx context.Context, eventID uint) error {
	message := EventCancelledMessage{
		EventID: eventID,
	}

	data, err := json.Marshal(message)
	if err != nil {
		p.logger.Error("failed to marshal event cancelled message",
			"error", err,
			"event_id", eventID)
		return err
	}

	kafkaMessage := kafka.Message{
		Topic: eventCancelled,
		Key:   []byte(fmt.Sprintf("%d", eventID)),
		Value: data,
		Time:  time.Now(),
	}

	// Retry logic
	maxRetries := 3
	for attempt := 0; attempt < maxRetries; attempt++ {
		err = p.writer.WriteMessages(ctx, kafkaMessage)
		if err == nil {
			p.logger.Info("event cancelled message sent",
				"event_id", eventID,
				"topic", eventCancelled)
			return nil
		}

		if attempt < maxRetries-1 {
			backoff := time.Duration(1<<uint(attempt)) * time.Second
			p.logger.Warn("failed to send event cancelled message, retrying",
				"error", err,
				"event_id", eventID,
				"attempt", attempt+1,
				"backoff", backoff)
			time.Sleep(backoff)
		}
	}

	p.logger.Error("failed to send event cancelled message after retries",
		"error", err,
		"event_id", eventID,
		"max_retries", maxRetries)
	return err
}

func (p *Producer) Close() error {
	if p.writer != nil {
		return p.writer.Close()
	}
	return nil
}

func (p *Producer) SendEventReminder(ctx context.Context, eventID uint, eventTitle string, eventDate time.Time) error {
	message := EventReminderMessage{
		EventID:    eventID,
		EventTitle: eventTitle,
		EventDate:  eventDate,
	}

	data, err := json.Marshal(message)
	if err != nil {
		p.logger.Error("failed to marshal event reminder message", "error", err, "event_id", eventID)
		return err
	}

	kafkaMessage := kafka.Message{
		Topic: eventReminder,
		Key:   []byte(fmt.Sprintf("%d", eventID)),
		Value: data,
		Time:  time.Now(),
	}

	// Retry logic
	maxRetries := 3
	for attempt := 0; attempt < maxRetries; attempt++ {
		err = p.writer.WriteMessages(ctx, kafkaMessage)
		if err == nil {
			p.logger.Info("event reminder message sent", "event_id", eventID, "topic", eventReminder)
			return nil
		}

		if attempt < maxRetries-1 {
			backoff := time.Duration(1<<uint(attempt)) * time.Second
			p.logger.Warn("failed to send event reminder message, retrying",
				"error", err,
				"event_id", eventID,
				"attempt", attempt+1,
				"backoff", backoff)
			time.Sleep(backoff)
		}
	}

	p.logger.Error("failed to send event reminder message after retries",
		"error", err,
		"event_id", eventID,
		"max_retries", maxRetries)
	return err
}
