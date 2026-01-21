package kafka

import (
	"context"
	"encoding/json"
	kafka "ticket-service/internal/kafka/events"

	kafka_go "github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka_go.Writer
}

func NewProducer(brokers []string) *Producer {
	return &Producer{
		writer: &kafka_go.Writer{
			Addr:     kafka_go.TCP(brokers...),
			Balancer: &kafka_go.LeastBytes{},
		},
	}
}

func (p *Producer) PublishTicketPurchased(
	ctx context.Context,
	event kafka.TicketPurchasedEvent,
) error {
	return p.writer.WriteMessages(ctx, kafka_go.Message{
		Topic: TopicTicketPurchased,
		Value: mustJSON(event),
	})
}

func (p *Producer) PublishTicketCheckin(
	ctx context.Context,
	event kafka.TicketCheckinEvent,
) error {
	return p.writer.WriteMessages(ctx, kafka_go.Message{
		Topic: TopicTicketCheckin,
		Value: mustJSON(event),
	})
}

func (p *Producer) Close() error {
	return p.writer.Close()
}

func mustJSON(v any) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}
