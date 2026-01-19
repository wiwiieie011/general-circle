package config

import (
	"os"
	"strings"
)

func KafkaBrokers() []string {
	brokers := os.Getenv("KAFKA_BROKER")
	if brokers == "" {
		return []string{"localhost:9092"}
	}
	return strings.Split(brokers, ",")
}
