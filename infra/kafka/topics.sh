#!/usr/bin/env bash
set -e

/opt/kafka/bin/kafka-topics.sh \
  --bootstrap-server localhost:9092 \
  --create --if-not-exists \
  --topic ticket.purchased \
  --partitions 3 \
  --replication-factor 1