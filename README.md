# 🎪 General Circle

## Платформа управления мероприятиями и билетами

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Docker](https://img.shields.io/badge/Docker-Compose-2496ED?style=flat&logo=docker)](https://docker.com)
[![Kafka](https://img.shields.io/badge/Apache-Kafka-231F20?style=flat&logo=apachekafka)](https://kafka.apache.org)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-14+-336791?style=flat&logo=postgresql)](https://www.postgresql.org)

## 📝 Описание

**General Circle** — микросервисная система для организации и управления мероприятиями. Платформа позволяет пользователям создавать события, покупать билеты и получать уведомления в реальном времени. Архитектура построена на Go с использованием Kafka для асинхронного взаимодействия между сервисами.

## 🚀 Технологии

- **Backend:** Go, Gin Framework
- **Messaging:** Apache Kafka
- **Database:** PostgreSQL, Redis
- **Gateway:** API Gateway с JWT middleware
- **Containerization:** Docker & Docker Compose
- **Authentication:** JWT (Bearer tokens)

## 🏗️ Микросервисы

- **Event Service** — управление мероприятиями и категориями
- **Ticket Service** — система продажи и управления билетами
- **User Service** — профили пользователей и аутентификация
- **Notification Service** — отправка уведомлений через Kafka
- **Gateway** — единая точка входа с JWT middleware

## 👥 Contributors

- 👨‍💻 - dzhambazbiev-ux
- 👨‍💻 - makaziev 
- 👨‍💻 - Strannik-chr
- 👨‍💻 - wiwiieie011

## 💡 Architecture Notes

- **Асинхронная передача данных** через Kafka для слабой связанности сервисов
- **JWT-токены** для безопасной аутентификации через Gateway
- **Repository pattern** для работы с БД
- **Service layer** для бизнес-логики
