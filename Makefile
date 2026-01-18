.PHONY: kafka-topics

# Запускает топики внутри kafka контейнера
kafka-topics:
	docker compose exec kafka /opt/kafka/topics.sh
