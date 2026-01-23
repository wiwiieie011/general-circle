#!/bin/bash
set -e

echo "Waiting for Kafka to be ready..."
sleep 10

KAFKA_HOME="/opt/kafka"
BROKER="localhost:9092"

echo "Creating topics..."

# Создание системного топика для consumer groups
$KAFKA_HOME/bin/kafka-topics.sh --bootstrap-server $BROKER --create --if-not-exists \
  --topic __consumer_offsets \
  --partitions 1 \
  --replication-factor 1 \
  --config cleanup.policy=compact || true

# Топики для event-service
$KAFKA_HOME/bin/kafka-topics.sh --bootstrap-server $BROKER --create --if-not-exists \
  --topic event.cancelled \
  --partitions 1 \
  --replication-factor 1 || true

$KAFKA_HOME/bin/kafka-topics.sh --bootstrap-server $BROKER --create --if-not-exists \
  --topic event.reminder \
  --partitions 1 \
  --replication-factor 1 || true

# Топики для ticket-service
$KAFKA_HOME/bin/kafka-topics.sh --bootstrap-server $BROKER --create --if-not-exists \
  --topic ticket.purchased \
  --partitions 1 \
  --replication-factor 1 || true

$KAFKA_HOME/bin/kafka-topics.sh --bootstrap-server $BROKER --create --if-not-exists \
  --topic ticket.checkin \
  --partitions 1 \
  --replication-factor 1 || true

echo "Topics created successfully!"
$KAFKA_HOME/bin/kafka-topics.sh --bootstrap-server $BROKER --list
