# Бизнес-флоу платформы "Сбор"

Документ описывает основные бизнес-процессы системы, их последовательность
и взаимодействие между сервисами.

---

## 1. Регистрация и авторизация пользователя

**Участники:** Client → Gateway → User Service

### Шаги
1. Пользователь отправляет `POST /api/auth/register`
2. Gateway проксирует запрос в User Service
3. User Service:
   - создаёт пользователя с ролью `user`
   - сохраняет в своей БД
4. Пользователь логинится через `POST /api/auth/login`
5. User Service:
   - валидирует креды
   - выдаёт `JWT access + refresh`
6. Gateway возвращает токены клиенту

**Результат:** пользователь авторизован

---

## 2. Получение профиля пользователя

**Участники:** Client → Gateway → User Service

### Шаги
1. Клиент отправляет `GET /api/users/me` с JWT
2. Gateway:
   - валидирует JWT
   - проксирует запрос
3. User Service возвращает данные пользователя

---

## 3. Становление организатором

**Участники:** Client → Gateway → User Service

### Шаги
1. Пользователь отправляет `POST /api/users/me/become-organizer`
2. Gateway валидирует JWT
3. User Service:
   - меняет роль `user → organizer`
   - сохраняет изменения

**Результат:** пользователь может создавать мероприятия

---

## 4. Создание мероприятия (draft)

**Участники:** Client → Gateway → Event Service

### Шаги
1. Организатор отправляет `POST /api/events`
2. Gateway:
   - валидирует JWT
   - проксирует запрос
3. Event Service:
   - создаёт Event в статусе `draft`
   - привязывает `organizer_id`
   - сохраняет в БД

**Результат:** мероприятие существует, но билеты продавать нельзя

---

## 5. Редактирование мероприятия

**Участники:** Client → Gateway → Event Service

### Шаги
1. Организатор отправляет `PUT /api/events/:id`
2. Event Service:
   - проверяет, что статус `draft`
   - обновляет данные

---

## 6. Публикация мероприятия

**Участники:** Client → Gateway → Event Service

### Шаги
1. Организатор отправляет `POST /api/events/:id/publish`
2. Event Service:
   - проверяет статус `draft`
   - переводит Event в `published`

**Результат:** разрешена продажа билетов

---

## 7. Создание типов билетов

**Участники:** Client → Gateway → Ticket Service → Event Service

### Шаги
1. Организатор отправляет `POST /api/events/:id/ticket-types`
2. Ticket Service:
   - по HTTP запрашивает Event Service
   - проверяет, что Event `published`
3. Ticket Service:
   - создаёт TicketType:
     - type (standard / vip / early_bird)
     - quantity
     - sales_start / sales_end
     - sold = 0

---

## 8. Покупка билета

**Участники:** Client → Gateway → Ticket Service → Event Service → Kafka

### Шаги
1. Пользователь отправляет `POST /api/events/:id/tickets`
2. Ticket Service:
   - проверяет Event (HTTP в Event Service)
   - проверяет период продаж
   - проверяет `sold < quantity`
3. Ticket Service:
   - увеличивает `sold`
   - создаёт Ticket:
     - user_id
     - event_id
     - ticket_type_id
     - status = active
     - уникальный `code` (для QR)
4. Ticket Service публикует событие:
   - `ticket.purchased` → Kafka

**Результат:** билет куплен

---

## 9. Уведомление о покупке билета

**Участники:** Kafka → Notification Service

### Шаги
1. Notification Service получает `ticket.purchased`
2. Создаёт Notification:
   - тип: "Билет куплен"
   - пользователь: `user_id`
3. Сохраняет в своей БД

---

## 10. Просмотр билетов пользователем

**Участники:** Client → Gateway → Ticket Service

### Шаги
1. Пользователь запрашивает `GET /api/tickets`
2. Ticket Service:
   - возвращает список билетов пользователя
   - с текущими статусами

---

## 11. Проверка билета (validate)

**Участники:** Scanner App → Gateway → Ticket Service

### Шаги
1. Приложение сканирует QR → получает `code`
2. Отправляет `POST /api/tickets/:code/validate`
3. Ticket Service:
   - проверяет существование билета
   - проверяет, что статус `active`
4. Возвращает результат (валиден / нет)

---

## 12. Использование билета (check-in)

**Участники:** Scanner App → Gateway → Ticket Service → Kafka

### Шаги
1. Отправляется `POST /api/tickets/:code/checkin`
2. Ticket Service:
   - проверяет, что билет `active`
   - меняет статус на `used`
3. Публикует событие:
   - `ticket.checkin` → Kafka

---

## 13. Отмена мероприятия

**Участники:** Client → Gateway → Event Service → Kafka

### Шаги
1. Организатор отправляет `POST /api/events/:id/cancel`
2. Event Service:
   - меняет статус на `cancelled`
3. Публикует событие:
   - `event.cancelled`

---

## 14. Уведомление об отмене мероприятия

**Участники:** Kafka → Notification Service → Ticket Service

### Шаги
1. Notification Service получает `event.cancelled`
2. Запрашивает Ticket Service:
   - список всех билетов по event_id
3. Для каждого пользователя:
   - создаёт Notification "Мероприятие отменено"

---

## 15. Напоминание о мероприятии

**Участники:** Event Service → Kafka → Notification Service

### Шаги
1. За день до события Event Service публикует:
   - `event.reminder`
2. Notification Service:
   - создаёт уведомления всем владельцам билетов

---

## Общая цепочка (коротко)

Client  
→ Gateway (JWT)  
→ Service (бизнес-логика)  
→ PostgreSQL (своя БД)  
→ Kafka (события)  
→ Notification Service  

---

## Важно

- JWT проверяется **только в Gateway**
- Сервисы доверяют Gateway
- Каждый сервис **НЕ ЛЕЗЕТ в чужую БД**
- События — только через Kafka