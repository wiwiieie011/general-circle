package config

import (
	"os"
)

type Config struct {
	App   AppConfig
	Kafka KafkaConfig
}

type AppConfig struct {
	Env string
}

type KafkaConfig struct {
	Brokers              []string
	TopicTicketPurchased string
}

func LoadConfig() *Config {
	return &Config{
		App: AppConfig{
			Env: os.Getenv("APP_ENV"),
		},
		Kafka: KafkaConfig{
			Brokers: []string{
				os.Getenv("KAFKA_BROKER"),
			},
			TopicTicketPurchased: os.Getenv("KAFKA_TOPIC_TICKET_PURCHASED"),
		},
	}
}
