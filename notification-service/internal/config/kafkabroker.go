package config

import "os"

func KafkaBrokers() []string {
	return []string{os.Getenv("KAFKA_BROKERS")}
}