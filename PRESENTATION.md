<div align="center">

# 🎯 Платформа "Сбор"

### Микросервисная система для организации мероприятий и продажи билетов

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)](https://go.dev/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-316192?style=for-the-badge&logo=postgresql)](https://www.postgresql.org/)
[![Kafka](https://img.shields.io/badge/Apache%20Kafka-3.7-231F20?style=for-the-badge&logo=apache-kafka)](https://kafka.apache.org/)
[![Docker](https://img.shields.io/badge/Docker-Compose-2496ED?style=for-the-badge&logo=docker)](https://www.docker.com/)

---

</div>

## 📑 Содержание

- [О проекте](#-о-проекте)
- [Архитектура системы](#️-архитектура-системы)
- [Компоненты системы](#-компоненты-системы)
- [Бизнес-процессы](#-бизнес-процессы)
- [Технологический стек](#️-технологический-стек)
- [API Endpoints](#-api-endpoints)
- [Безопасность](#-безопасность)
- [Запуск проекта](#-запуск-проекта)

---

<div align="center">

## 📋 О проекте

</div>

**"Сбор"** — современная микросервисная платформа для полного цикла управления мероприятиями: от создания до продажи билетов и проверки на входе.

### ✨ Ключевые возможности

<table>
<tr>
<td width="50%">

#### 🔐 Аутентификация
- Регистрация и вход пользователей
- JWT токены (access + refresh)
- Управление ролями (user/organizer)

#### 🎫 Управление билетами
- Создание типов билетов
- Продажа билетов онлайн
- QR-коды для входа
- Валидация и check-in

</td>
<td width="50%">

#### 📅 Мероприятия
- Создание и редактирование
- Публикация и отмена
- Расписание и категории
- Автоматические напоминания

#### 📱 Уведомления
- Push-уведомления
- Email-рассылки
- Настройки пользователя
- История уведомлений

</td>
</tr>
</table>

---

<div align="center">

## 🏗️ Архитектура системы

</div>

### Общая архитектура

```mermaid
graph TB
    subgraph Client[" "]
        Web[🌐 Web App]
        Mobile[📱 Mobile App]
    end
    
    subgraph Gateway["🔐 API Gateway :8000"]
        JWT[JWT Middleware]
        Router[Router]
    end
    
    subgraph Services["⚙️ Микросервисы"]
        direction TB
        US[👤 User Service<br/>:8081<br/>Auth & Users]
        ES[📅 Event Service<br/>:8083<br/>Events & Schedule]
        TS[🎫 Ticket Service<br/>:8082<br/>Tickets & Types]
        NS[📬 Notification Service<br/>:8084<br/>Notifications]
    end
    
    subgraph Data["💾 Базы данных"]
        PG1[(PostgreSQL<br/>User DB)]
        PG2[(PostgreSQL<br/>Event DB)]
        PG3[(PostgreSQL<br/>Ticket DB)]
        PG4[(PostgreSQL<br/>Notification DB)]
        Redis[(Redis<br/>Cache)]
    end
    
    subgraph Messaging["📨 Message Broker"]
        Kafka[Apache Kafka<br/>:9092<br/>Event Streaming]
    end
    
    Web --> Gateway
    Mobile --> Gateway
    Gateway --> JWT
    JWT --> Router
    
    Router -->|/api/auth/*| US
    Router -->|/api/users/*| US
    Router -->|/api/events/*| ES
    Router -->|/api/tickets/*| TS
    Router -->|/api/notifications/*| NS
    
    US --> PG1
    ES --> PG2
    TS --> PG3
    NS --> PG4
    NS --> Redis
    
    ES -.->|HTTP| TS
    ES -->|Publish| Kafka
    TS -->|Publish| Kafka
    Kafka -->|Consume| NS
    
    style Gateway fill:#4A90E2,stroke:#2E5C8A,color:#fff
    style US fill:#50C878,stroke:#2D8659,color:#fff
    style ES fill:#FF6B6B,stroke:#C92A2A,color:#fff
    style TS fill:#FFD93D,stroke:#F59F00,color:#fff
    style NS fill:#9B59B6,stroke:#6C3483,color:#fff
    style Kafka fill:#000,stroke:#FF6B00,color:#fff
```

### Принципы архитектуры

<table>
<tr>
<td>

**🔒 Изоляция сервисов**
- Каждый сервис имеет свою БД
- Нет прямого доступа к чужим данным
- Независимое масштабирование

</td>
<td>

**🔄 Асинхронная коммуникация**
- Kafka для событий
- HTTP для синхронных запросов
- Event-driven архитектура

</td>
</tr>
<tr>
<td>

**🛡️ Безопасность**
- Централизованная JWT валидация
- Gateway как единая точка входа
- Изоляция сервисов

</td>
<td>

**📈 Масштабируемость**
- Горизонтальное масштабирование
- Независимое развертывание
- Микросервисная архитектура

</td>
</tr>
</table>

---

<div align="center">

## 🧩 Компоненты системы

</div>

### 🔐 API Gateway (порт 8000)

<div align="center">

| Функция | Описание |
|---------|----------|
| 🚪 **Единая точка входа** | Все запросы проходят через Gateway |
| 🔑 **JWT валидация** | Проверка токенов для защищенных эндпоинтов |
| 🧭 **Маршрутизация** | Проксирование запросов к нужным сервисам |
| 🛡️ **Безопасность** | Защита от несанкционированного доступа |

</div>

### 👤 User Service (порт 8081)

```mermaid
graph LR
    A[Client] -->|Register/Login| B[User Service]
    B -->|Create User| C[(User DB)]
    B -->|Generate JWT| D[Access Token]
    B -->|Generate JWT| E[Refresh Token]
    D --> A
    E --> A
    
    style B fill:#50C878,stroke:#2D8659,color:#fff
    style C fill:#4A90E2,stroke:#2E5C8A,color:#fff
```

**Функции:**
- ✅ Регистрация новых пользователей
- ✅ Аутентификация (login)
- ✅ Управление профилями
- ✅ Изменение роли (user → organizer)
- ✅ Выдача JWT токенов

### 📅 Event Service (порт 8083)

```mermaid
graph TB
    A[Organizer] -->|Create Event| B[Event Service]
    B -->|Save| C[(Event DB)]
    B -->|Publish| D[Published Event]
    D -->|Schedule| E[Cron Job]
    E -->|Reminder| F[Kafka]
    B -->|Cancel| F
    
    style B fill:#FF6B6B,stroke:#C92A2A,color:#fff
    style C fill:#4A90E2,stroke:#2E5C8A,color:#fff
    style F fill:#000,stroke:#FF6B00,color:#fff
```

**Функции:**
- ✅ Создание мероприятий (draft → published)
- ✅ Управление категориями
- ✅ Расписание мероприятий
- ✅ Автоматические напоминания (cron)
- ✅ Публикация событий в Kafka

### 🎫 Ticket Service (порт 8082)

```mermaid
graph TB
    A[User] -->|Buy Ticket| B[Ticket Service]
    B -->|Check Event| C[Event Service]
    C -->|Event Status| B
    B -->|Validate| D{Available?}
    D -->|Yes| E[Create Ticket]
    D -->|No| F[Error]
    E -->|Generate QR| G[Unique Code]
    E -->|Save| H[(Ticket DB)]
    E -->|Publish| I[Kafka]
    
    style B fill:#FFD93D,stroke:#F59F00,color:#fff
    style C fill:#FF6B6B,stroke:#C92A2A,color:#fff
    style H fill:#4A90E2,stroke:#2E5C8A,color:#fff
    style I fill:#000,stroke:#FF6B00,color:#fff
```

**Функции:**
- ✅ Создание типов билетов
- ✅ Продажа билетов
- ✅ Генерация QR-кодов
- ✅ Валидация билетов
- ✅ Check-in на мероприятии

### 📬 Notification Service (порт 8084)

```mermaid
graph LR
    A[Kafka] -->|Events| B[Notification Service]
    B -->|Create| C[(Notification DB)]
    B -->|Cache| D[(Redis)]
    B -->|Send| E[User]
    
    style B fill:#9B59B6,stroke:#6C3483,color:#fff
    style C fill:#4A90E2,stroke:#2E5C8A,color:#fff
    style D fill:#DC2626,stroke:#991B1B,color:#fff
    style A fill:#000,stroke:#FF6B00,color:#fff
```

**Функции:**
- ✅ Подписка на Kafka события
- ✅ Создание уведомлений
- ✅ Кэширование в Redis
- ✅ Управление настройками
- ✅ История уведомлений

---

<div align="center">

## 🔄 Бизнес-процессы

</div>

### 1. Жизненный цикл мероприятия

```mermaid
stateDiagram-v2
    [*] --> Draft: Организатор создает
    Draft --> Draft: Редактирование данных
    Draft --> Published: Публикация
    Published --> Published: Продажа билетов
    Published --> Cancelled: Отмена организатором
    Published --> [*]: Завершение мероприятия
    Cancelled --> [*]: Мероприятие отменено
    
    note right of Draft
        Статус: draft
        Билеты не продаются
    end note
    
    note right of Published
        Статус: published
        Билеты доступны
    end note
    
    note right of Cancelled
        Статус: cancelled
        Все билеты аннулированы
    end note
```

### 2. Жизненный цикл билета

```mermaid
stateDiagram-v2
    [*] --> Active: Покупка билета
    Active --> Active: Валидация QR-кода
    Active --> Used: Check-in на входе
    Active --> Cancelled: Отмена мероприятия
    Used --> [*]: Билет использован
    Cancelled --> [*]: Билет аннулирован
    
    note right of Active
        Статус: active
        Можно использовать
    end note
    
    note right of Used
        Статус: used
        Билет использован
    end note
```

### 3. Процесс покупки билета (детальный)

```mermaid
flowchart TD
    Start([👤 Пользователь<br/>хочет купить билет]) --> Auth{🔐 Проверка<br/>JWT токена}
    Auth -->|Невалидный| RejectAuth[❌ Ошибка:<br/>Не авторизован]
    Auth -->|Валидный| CheckEvent{📅 Проверить<br/>статус мероприятия}
    
    CheckEvent -->|draft| Reject1[❌ Ошибка:<br/>Мероприятие не опубликовано]
    CheckEvent -->|cancelled| Reject2[❌ Ошибка:<br/>Мероприятие отменено]
    CheckEvent -->|published| CheckSales{⏰ Проверить<br/>период продаж}
    
    CheckSales -->|Не начались| Reject3[❌ Ошибка:<br/>Продажи еще не начались]
    CheckSales -->|Закончились| Reject4[❌ Ошибка:<br/>Продажи закончились]
    CheckSales -->|Активны| CheckQuantity{🎫 Проверить<br/>количество билетов}
    
    CheckQuantity -->|Нет мест| Reject5[❌ Ошибка:<br/>Билеты распроданы]
    CheckQuantity -->|Есть места| CreateTicket[✅ Создать билет]
    
    CreateTicket --> GenerateQR[🔲 Генерация<br/>уникального QR-кода]
    GenerateQR --> IncrementSold[📊 Увеличить<br/>счетчик sold]
    IncrementSold --> SaveDB[💾 Сохранить<br/>в БД]
    SaveDB --> PublishEvent[📨 Опубликовать<br/>ticket.purchased]
    PublishEvent --> SendNotification[📬 Создать<br/>уведомление]
    SendNotification --> Success([🎉 Билет<br/>успешно куплен!])
    
    RejectAuth --> End([Конец])
    Reject1 --> End
    Reject2 --> End
    Reject3 --> End
    Reject4 --> End
    Reject5 --> End
    Success --> End
    
    style Start fill:#E8F5E9,stroke:#4CAF50
    style Success fill:#C8E6C9,stroke:#2E7D32
    style RejectAuth fill:#FFEBEE,stroke:#F44336
    style Reject1 fill:#FFEBEE,stroke:#F44336
    style Reject2 fill:#FFEBEE,stroke:#F44336
    style Reject3 fill:#FFEBEE,stroke:#F44336
    style Reject4 fill:#FFEBEE,stroke:#F44336
    style Reject5 fill:#FFEBEE,stroke:#F44336
```

### 4. Полный цикл: от регистрации до использования билета

```mermaid
sequenceDiagram
    autonumber
    participant U as 👤 User
    participant G as 🔐 Gateway
    participant US as 👤 User Service
    participant O as 🎭 Organizer
    participant ES as 📅 Event Service
    participant TS as 🎫 Ticket Service
    participant K as 📨 Kafka
    participant NS as 📬 Notification Service
    
    rect rgb(200, 230, 201)
        Note over U,US: 🔐 Регистрация и авторизация
        U->>G: POST /api/auth/register
        G->>US: Proxy request
        US->>US: Create user (role: user)
        US-->>G: User created
        G-->>U: ✅ Success
        
        U->>G: POST /api/auth/login
        G->>US: Proxy request
        US->>US: Validate credentials
        US-->>G: JWT tokens
        G-->>U: 🔑 Access + Refresh tokens
    end
    
    rect rgb(255, 235, 238)
        Note over O,ES: 📅 Создание мероприятия
        O->>G: POST /api/events (JWT)
        G->>G: ✅ Validate JWT
        G->>ES: Proxy request
        ES->>ES: Create event (status: draft)
        ES-->>G: Event created
        G-->>O: 📅 Event (draft)
        
        O->>G: POST /api/events/:id/publish
        G->>ES: Proxy request
        ES->>ES: Change status → published
        ES-->>G: ✅ Published
        G-->>O: 📅 Event (published)
    end
    
    rect rgb(255, 243, 224)
        Note over O,TS: 🎫 Создание типов билетов
        O->>G: POST /api/events/:id/ticket-types
        G->>TS: Proxy request
        TS->>ES: HTTP: Check event status
        ES-->>TS: Event (published)
        TS->>TS: Create TicketType
        TS-->>G: ✅ TicketType created
        G-->>O: 🎫 TicketType
    end
    
    rect rgb(232, 245, 233)
        Note over U,NS: 💰 Покупка билета
        U->>G: POST /api/events/:id/tickets (JWT)
        G->>G: ✅ Validate JWT
        G->>TS: Proxy request
        TS->>ES: HTTP: Get event details
        ES-->>TS: Event details
        TS->>TS: Validate & Create ticket
        TS->>TS: Generate QR code
        TS->>K: 📨 Publish ticket.purchased
        TS-->>G: ✅ Ticket created
        G-->>U: 🎫 Ticket with QR
        
        K->>NS: ticket.purchased event
        NS->>NS: Create notification
        NS->>NS: Store in DB
        Note over NS: 📬 Notification ready
    end
    
    rect rgb(243, 229, 245)
        Note over U,TS: ✅ Проверка билета
        U->>G: POST /api/tickets/:code/validate
        G->>TS: Proxy request
        TS->>TS: Validate ticket
        TS-->>G: ✅ Valid / ❌ Invalid
        G-->>U: Validation result
        
        U->>G: POST /api/tickets/:code/checkin
        G->>TS: Proxy request
        TS->>TS: Change status → used
        TS->>K: 📨 Publish ticket.checkin
        TS-->>G: ✅ Check-in successful
        G-->>U: 🎉 Welcome!
    end
```

### 5. Система уведомлений (Event-Driven)

```mermaid
graph TB
    subgraph Producers["📤 Event Producers"]
        ES[📅 Event Service]
        TS[🎫 Ticket Service]
    end
    
    subgraph Broker["📨 Kafka Message Broker"]
        K[Apache Kafka<br/>Topics]
    end
    
    subgraph Consumer["📥 Notification Service"]
        NS[Notification Service]
        Handler[Event Handler]
    end
    
    subgraph Storage["💾 Storage"]
        DB[(PostgreSQL<br/>Notifications)]
        Redis[(Redis<br/>Cache)]
    end
    
    subgraph Users["👥 Users"]
        U1[User 1]
        U2[User 2]
        U3[User N...]
    end
    
    ES -->|event.cancelled| K
    ES -->|event.reminder| K
    TS -->|ticket.purchased| K
    TS -->|ticket.checkin| K
    
    K -->|Consume| NS
    NS --> Handler
    Handler -->|Create| DB
    Handler -->|Cache| Redis
    Handler -->|Notify| U1
    Handler -->|Notify| U2
    Handler -->|Notify| U3
    
    style ES fill:#FF6B6B,stroke:#C92A2A,color:#fff
    style TS fill:#FFD93D,stroke:#F59F00,color:#fff
    style K fill:#000,stroke:#FF6B00,color:#fff
    style NS fill:#9B59B6,stroke:#6C3483,color:#fff
    style DB fill:#4A90E2,stroke:#2E5C8A,color:#fff
    style Redis fill:#DC2626,stroke:#991B1B,color:#fff
```

---

<div align="center">

## 🛠️ Технологический стек

</div>

### Backend

<table>
<tr>
<td width="33%">

#### 🟢 Go (Golang)
- Версия: 1.21+
- Фреймворк: Gin
- ORM: GORM
- Логирование: slog

</td>
<td width="33%">

#### 🐘 PostgreSQL
- Версия: 16
- Database per Service
- Миграции: GORM AutoMigrate
- Изоляция данных

</td>
<td width="33%">

#### 🔴 Redis
- Кэширование уведомлений
- Быстрый доступ к данным
- Session storage

</td>
</tr>
<tr>
<td>

#### 📨 Apache Kafka
- Версия: 3.7.0
- Event-driven архитектура
- Асинхронная коммуникация
- Надежная доставка

</td>
<td>

#### 🔐 JWT
- Access tokens
- Refresh tokens
- Безопасная аутентификация
- Stateless авторизация

</td>
<td>

#### 🐳 Docker
- Контейнеризация
- Docker Compose
- Изолированная среда
- Легкое развертывание

</td>
</tr>
</table>

### Архитектурные паттерны

<div align="center">

| Паттерн | Описание | Применение |
|---------|----------|------------|
| 🏗️ **Microservices** | Разделение на независимые сервисы | Все сервисы изолированы |
| 🚪 **API Gateway** | Единая точка входа | Gateway :8000 |
| 📨 **Event-Driven** | Асинхронная коммуникация через события | Kafka events |
| 💾 **Database per Service** | Каждый сервис имеет свою БД | Изоляция данных |
| 🔄 **CQRS** | Разделение чтения и записи | Оптимизация запросов |
| ⏰ **Cron Jobs** | Планировщик задач | Напоминания о мероприятиях |

</div>

---

<div align="center">

## 📡 API Endpoints

</div>

### 🔐 Gateway (порт 8000)

<div align="center">

| Метод | Путь | Назначение | Сервис |
|-------|------|------------|--------|
| `POST` | `/api/auth/*` | Аутентификация | User Service |
| `GET` | `/api/users/*` | Управление пользователями | User Service |
| `GET/POST/PUT` | `/api/events/*` | Управление мероприятиями | Event Service |
| `GET/POST` | `/api/tickets/*` | Управление билетами | Ticket Service |
| `GET` | `/api/notifications/*` | Уведомления | Notification Service |

</div>

### 👤 User Service (порт 8081)

<table>
<tr>
<td>

**Аутентификация**
- `POST /auth/register` - Регистрация
- `POST /auth/login` - Вход в систему
- `POST /auth/refresh` - Обновление токена

</td>
<td>

**Профиль**
- `GET /users/me` - Получить профиль
- `PUT /users/me` - Обновить профиль
- `POST /users/me/become-organizer` - Стать организатором

</td>
</tr>
</table>

### 📅 Event Service (порт 8083)

<table>
<tr>
<td>

**Мероприятия**
- `POST /events` - Создать мероприятие
- `GET /events` - Список мероприятий
- `GET /events/:id` - Детали мероприятия
- `PUT /events/:id` - Обновить мероприятие

</td>
<td>

**Управление**
- `POST /events/:id/publish` - Опубликовать
- `POST /events/:id/cancel` - Отменить
- `GET /events/:id/schedule` - Расписание

</td>
</tr>
</table>

### 🎫 Ticket Service (порт 8082)

<table>
<tr>
<td>

**Типы билетов**
- `POST /events/:id/ticket-types` - Создать тип
- `GET /events/:id/ticket-types` - Список типов

</td>
<td>

**Билеты**
- `POST /events/:id/tickets` - Купить билет
- `GET /tickets` - Мои билеты
- `POST /tickets/:code/validate` - Проверить
- `POST /tickets/:code/checkin` - Использовать

</td>
</tr>
</table>

### 📬 Notification Service (порт 8084)

<table>
<tr>
<td>

**Уведомления**
- `GET /notifications` - Получить список
- `GET /notifications/:id` - Детали
- `PUT /notifications/:id/read` - Отметить прочитанным

</td>
<td>

**Настройки**
- `GET /notifications/preferences` - Настройки
- `PUT /notifications/preferences` - Обновить

</td>
</tr>
</table>

---

<div align="center">

## 🔐 Безопасность

</div>

### Архитектура безопасности

```mermaid
graph TB
    Client[👤 Client] -->|Request| Gateway[🔐 Gateway]
    Gateway -->|Validate JWT| JWT{🔑 JWT Valid?}
    JWT -->|No| Reject[❌ 401 Unauthorized]
    JWT -->|Yes| Extract[📋 Extract User Info]
    Extract -->|Proxy| Service[⚙️ Microservice]
    Service -->|Trust Gateway| Process[✅ Process Request]
    
    style Gateway fill:#4A90E2,stroke:#2E5C8A,color:#fff
    style JWT fill:#50C878,stroke:#2D8659,color:#fff
    style Reject fill:#FF6B6B,stroke:#C92A2A,color:#fff
```

### Принципы безопасности

<div align="center">

| ✅ Принцип | 📝 Описание |
|-----------|-------------|
| **Централизованная валидация** | JWT проверяется только в Gateway |
| **Доверие между сервисами** | Сервисы доверяют Gateway |
| **Изоляция данных** | Каждый сервис имеет свою БД |
| **Нет прямого доступа** | Сервисы не лезут в чужие БД |
| **Безопасная коммуникация** | HTTP для синхронных, Kafka для асинхронных |
| **Роли и права** | user, organizer с разными правами |

</div>

### JWT Токены

<table>
<tr>
<td>

**Access Token**
- Короткий срок жизни
- Используется для запросов
- Содержит user_id, role
- Валидируется в Gateway

</td>
<td>

**Refresh Token**
- Долгий срок жизни
- Используется для обновления
- Хранится безопасно
- Обновляет access token

</td>
</tr>
</table>

---

<div align="center">

## 📈 Kafka Events

</div>

### События системы

<div align="center">

| 📨 Событие | 📤 Источник | 📥 Получатель | 📝 Описание |
|-----------|------------|--------------|-------------|
| `ticket.purchased` | 🎫 Ticket Service | 📬 Notification Service | Билет успешно куплен |
| `ticket.checkin` | 🎫 Ticket Service | 📬 Notification Service | Билет использован на входе |
| `event.cancelled` | 📅 Event Service | 📬 Notification Service | Мероприятие отменено |
| `event.reminder` | 📅 Event Service | 📬 Notification Service | Напоминание о мероприятии |

</div>

### Пример события

```json
{
  "event": "ticket.purchased",
  "timestamp": "2026-01-18T11:00:00Z",
  "data": {
    "ticket_id": "123e4567-e89b-12d3-a456-426614174000",
    "event_id": 42,
    "ticket_type_id": 3,
    "user_id": 7,
    "code": "ABCD1234",
    "status": "active",
    "purchased_at": "2026-01-18T11:00:00Z"
  }
}
```

---

<div align="center">

## 🚀 Запуск проекта

</div>

### Требования

<div align="center">

| Компонент | Версия | Назначение |
|-----------|--------|------------|
| 🐳 Docker | Latest | Контейнеризация |
| 🐳 Docker Compose | Latest | Оркестрация |
| 🟢 Go | 1.21+ | Backend разработка |
| 🐘 PostgreSQL | 16 | База данных |
| 📨 Kafka | 3.7.0 | Message broker |

</div>

### Быстрый старт

#### 1️⃣ Запуск инфраструктуры

```bash
# Запустить PostgreSQL, Kafka, Redis
docker-compose up -d

# Проверить статус
docker-compose ps
```

#### 2️⃣ Запуск сервисов

```bash
# Terminal 1: Gateway
cd gateway && go run main.go

# Terminal 2: User Service
cd user-service && go run cmd/app/main.go

# Terminal 3: Event Service
cd event-service && go run cmd/app/main.go

# Terminal 4: Ticket Service
cd ticket-service && go run cmd/app/main.go

# Terminal 5: Notification Service
cd notification-service && go run cmd/app/main.go
```

### Порты сервисов

<div align="center">

| Сервис | Порт | URL |
|--------|------|-----|
| 🔐 Gateway | 8000 | http://localhost:8000 |
| 👤 User Service | 8081 | http://localhost:8081 |
| 🎫 Ticket Service | 8082 | http://localhost:8082 |
| 📅 Event Service | 8083 | http://localhost:8083 |
| 📬 Notification Service | 8084 | http://localhost:8084 |
| 🐘 PostgreSQL | 5432 | localhost:5432 |
| 📨 Kafka | 9092 | localhost:9092 |
| 🎨 Kafka UI | 8090 | http://localhost:8090 |

</div>

---

<div align="center">

## 👥 Роли пользователей

</div>

### 👤 User (Пользователь)

<table>
<tr>
<td width="50%">

**Возможности:**
- ✅ Регистрация и авторизация
- ✅ Просмотр мероприятий
- ✅ Покупка билетов
- ✅ Просмотр своих билетов
- ✅ Получение уведомлений
- ✅ Валидация QR-кодов

</td>
<td width="50%">

**Ограничения:**
- ❌ Не может создавать мероприятия
- ❌ Не может управлять билетами
- ❌ Не может отменять события

</td>
</tr>
</table>

### 🎭 Organizer (Организатор)

<table>
<tr>
<td width="50%">

**Все возможности User +**
- ✅ Создание мероприятий
- ✅ Редактирование мероприятий
- ✅ Публикация/отмена событий
- ✅ Создание типов билетов
- ✅ Просмотр статистики

</td>
<td width="50%">

**Дополнительно:**
- ✅ Управление расписанием
- ✅ Управление категориями
- ✅ Полный контроль над событиями

</td>
</tr>
</table>

---

<div align="center">

## 🎯 Ключевые особенности

</div>

<table>
<tr>
<td width="33%">

### 📈 Масштабируемость
- Горизонтальное масштабирование
- Независимое масштабирование сервисов
- Load balancing готовность

</td>
<td width="33%">

### 🛡️ Отказоустойчивость
- Изоляция сервисов
- Нет каскадных сбоев
- Graceful degradation

</td>
<td width="33%">

### ⚡ Производительность
- Кэширование в Redis
- Асинхронная обработка
- Оптимизированные запросы

</td>
</tr>
<tr>
<td>

### 🔄 Гибкость
- Легко добавлять новые сервисы
- Независимое развертывание
- Модульная архитектура

</td>
<td>

### 🔐 Безопасность
- Централизованная аутентификация
- Изоляция данных
- JWT токены

</td>
<td>

### 📊 Мониторинг
- Structured logging
- Event tracking через Kafka
- Готовность к метрикам

</td>
</tr>
</table>

---

<div align="center">

## 📦 Структура проекта

</div>

```
general-circle/
├── 🔐 gateway/                    # API Gateway
│   ├── main.go
│   ├── proxy.go
│   ├── middleware/
│   │   └── jwt.go
│   └── jwtutil/
│       └── jwt.go
│
├── 👤 user-service/               # User Service
│   ├── cmd/app/main.go
│   ├── internal/
│   │   ├── config/
│   │   ├── models/
│   │   ├── repository/
│   │   ├── services/
│   │   └── transport/
│   └── go.mod
│
├── 📅 event-service/              # Event Service
│   ├── cmd/app/main.go
│   ├── internal/
│   │   ├── config/
│   │   ├── models/
│   │   ├── repository/
│   │   ├── services/
│   │   ├── transport/
│   │   └── kafka/
│   └── go.mod
│
├── 🎫 ticket-service/             # Ticket Service
│   ├── cmd/app/main.go
│   ├── internal/
│   │   ├── api/http/
│   │   ├── config/
│   │   ├── models/
│   │   ├── repository/
│   │   ├── services/
│   │   ├── transport/
│   │   └── kafka/
│   └── go.mod
│
├── 📬 notification-service/      # Notification Service
│   ├── cmd/app/main.go
│   ├── internal/
│   │   ├── config/
│   │   ├── models/
│   │   ├── repository/
│   │   ├── services/
│   │   ├── transport/
│   │   └── kafka/
│   └── go.mod
│
├── 🏗️ infra/                     # Инфраструктура
│   └── kafka/
│       ├── server.properties
│       └── topics.sh
│
├── 📚 docs/                       # Документация
│   └── api-contracts.md
│
├── 📄 BUSINESS_FLOWS.md           # Бизнес-процессы
├── 🐳 docker-compose.yaml         # Docker Compose
└── 📋 PRESENTATION.md             # Эта презентация
```

---

<div align="center">

## 📝 Дополнительная документация

</div>

<div align="center">

| 📄 Документ | 📝 Описание |
|------------|-------------|
| [BUSINESS_FLOWS.md](./BUSINESS_FLOWS.md) | Подробное описание всех бизнес-процессов |
| [API Контракты](./docs/api-contracts.md) | Детальные спецификации API |

</div>

---

<div align="center">

## 🎉 Заключение

**"Сбор"** — это современная, масштабируемая и надежная платформа для управления мероприятиями, построенная на принципах микросервисной архитектуры и event-driven подхода.

### Основные преимущества:
- 🚀 Быстрое развертывание
- 📈 Легкое масштабирование
- 🛡️ Высокая безопасность
- ⚡ Отличная производительность
- 🔄 Гибкая архитектура

---

*Документ создан для презентации проекта "Сбор"*  
*Версия: 1.0 | Дата: 2026*

</div>
