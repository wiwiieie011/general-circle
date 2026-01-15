# API Контракты

## Ticket Service (порт 8082)

### GET /ping
Запрос:
Ответ: { "status": "success" }

### GET /events/:id/ticket-types
Запрос:
{
  "type": "vip",
  "price": 5000,
  "quantity": 100,
  "sales_start": "2026-02-01T00:00:00Z",
  "sales_end": "2026-02-10T23:59:59Z"
}
Ответ:
{
  "id": 12,
  "event_id": 5,
  "type": "vip",
  "price": 5000,
  "quantity": 100,
  "sold": 23,
  "sales_start": "2026-02-01T00:00:00Z",
  "sales_end": "2026-02-10T23:59:59Z",
  "created_at": "2026-01-10T12:00:00Z",
  "updated_at": "2026-01-10T12:00:00Z"
}

## Kafka-события
