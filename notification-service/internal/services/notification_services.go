package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	// "encoding/json"
	"errors"
	// "fmt"
	"log/slog"
	"notification-service/internal/dto"
	"notification-service/internal/models"
	"notification-service/internal/repository"

	"github.com/redis/go-redis/v9"
	// "time"
	// "github.com/redis/go-redis/v9"
)

type NotificationService interface {
	CreateNotificationInternal(notification *models.Notification) error
	GetNotifications(ctx context.Context, userID uint, limit int, lastID uint) ([]models.Notification, error)
	CheckAll(userID uint) error
	CheckNotificationsByID(userID, id uint) error
	DeleteNotificationByID(userID, id uint) error
	GetNotificationPreferences(userID uint) (*models.NotificationPreference, error)
	Update(userID uint, req dto.UpdateNotificationPreferencesRequest) (*models.NotificationPreference, error)
	Count(userID uint) (int64, error)
}

type notificationService struct {
	notificationRepo repository.NotificationRepo
	redis            *redis.Client
	log              *slog.Logger
}

func NewNotificationService(notificationRepo repository.NotificationRepo, log *slog.Logger, redis *redis.Client) NotificationService {
	return &notificationService{
		notificationRepo: notificationRepo,
		redis:            redis,
		log:              log,
	}
}

func (s *notificationService) CreateNotificationInternal(notification *models.Notification) error {
	if notification.UserID == 0 {
		s.log.Warn("create notification failed: invalid user id")
		return errors.New("invalid user id")
	}

	notification.Read = false

	if err := s.notificationRepo.Create(notification); err != nil {
		s.log.Error(
			"failed to create notification",
			"error", err,
			"userID", notification.UserID,
		)
		return err
	}

	s.log.Info(
		"notification created",
		"userID", notification.UserID,
		"notificationID", notification.ID,
	)

	return nil
}

func (s *notificationService) GetNotifications(ctx context.Context, userID uint, limit int, lastID uint) ([]models.Notification, error) {

	if userID == 0 {
		s.log.Warn("get notifications unauthorized")
		return nil, dto.ErrUnauthorized
	}

	if lastID == 0 {
		// Если Redis не инициализирован (например, в unit-тестах) — обходим кэш
		if s.redis == nil {
			nots, err := s.notificationRepo.GetNotifications(userID, limit, 0)
			if err != nil {
				s.log.Error("failed to get notifications", "error", err, "userID", userID)
				return nil, err
			}
			s.log.Info("notifications fetched from db (redis disabled)", "userID", userID, "count", len(nots))
			return nots, nil
		}

		cacheKey := fmt.Sprintf("notifications:%d:first", userID)
		cached, err := s.redis.Get(ctx, cacheKey).Result()
		if err == nil && cached != "" {
			var nots []models.Notification
			if err := json.Unmarshal([]byte(cached), &nots); err == nil {
				s.log.Info("notifications fetched from cache", "userID", userID, "count", len(nots))
				return nots, nil
			}
		}
		nots, err := s.notificationRepo.GetNotifications(userID, limit, 0)
		if err != nil {
			s.log.Error("failed to get notifications", "error", err, "userID", userID)
			return nil, err
		}

		data, _ := json.Marshal(nots)
		_ = s.redis.Set(ctx, cacheKey, data, 5*time.Second).Err()

		s.log.Info("notifications fetched from db and cached", "userID", userID, "count", len(nots))
		return nots, nil
	}

	// 🔹 ВСЕ ОСТАЛЬНЫЕ СТРАНИЦЫ — ТОЛЬКО ИЗ БД (БЕЗ КЭША)
	nots, err := s.notificationRepo.GetNotifications(userID, limit, lastID)
	if err != nil {
		s.log.Error("failed to get notifications", "error", err, "userID", userID, "lastID", lastID)
		return nil, err
	}

	s.log.Info("notifications fetched from db (no cache)", "userID", userID, "count", len(nots))
	return nots, nil
}

func (s *notificationService) CheckAll(userID uint) error {
	if userID == 0 {
		s.log.Warn("check all notifications unauthorized")
		return dto.ErrUnauthorized
	}

	if err := s.notificationRepo.AllRead(userID); err != nil {
		s.log.Error(
			"failed to mark all notifications as read",
			"error", err,
			"userID", userID,
		)
		return err
	}

	s.log.Info(
		"all notifications marked as read",
		"userID", userID,
	)

	return nil
}

func (s *notificationService) CheckNotificationsByID(userID, id uint) error {
	if userID == 0 {
		s.log.Warn("mark notification read unauthorized")
		return dto.ErrUnauthorized
	}

	if id == 0 {
		s.log.Warn(
			"invalid notification id",
			"userID", userID,
		)
		return dto.ErrInvalidNotificationID
	}

	if err := s.notificationRepo.ReadNotificationsByID(userID, id); err != nil {
		s.log.Error(
			"failed to mark notification as read",
			"error", err,
			"userID", userID,
			"notificationID", id,
		)
		return err
	}

	s.log.Info(
		"notification marked as read",
		"userID", userID,
		"notificationID", id,
	)

	return nil
}

func (s *notificationService) DeleteNotificationByID(userID, id uint) error {
	if userID == 0 {
		s.log.Warn("delete notification unauthorized")
		return dto.ErrUnauthorized
	}

	if id == 0 {
		s.log.Warn(
			"invalid notification id for delete",
			"userID", userID,
		)
		return dto.ErrInvalidNotificationID
	}

	if err := s.notificationRepo.DeleteNotificationsByID(userID, id); err != nil {
		s.log.Error(
			"failed to delete notification",
			"error", err,
			"userID", userID,
			"notificationID", id,
		)
		return err
	}

	s.log.Info(
		"notification deleted",
		"userID", userID,
		"notificationID", id,
	)

	return nil
}

func (s *notificationService) GetNotificationPreferences(userID uint) (*models.NotificationPreference, error) {
	if userID == 0 {
		s.log.Warn("get notification preferences unauthorized")
		return nil, dto.ErrUnauthorized
	}

	val, err := s.notificationRepo.GetNotificationPreferences(userID)
	if err != nil {
		s.log.Error(
			"failed to get notification preferences",
			"error", err,
			"userID", userID,
		)
		return nil, dto.ErrPreferencesNotFound
	}

	s.log.Info(
		"notification preferences fetched",
		"userID", userID,
	)

	return val, nil
}

func (s *notificationService) Update(userID uint, req dto.UpdateNotificationPreferencesRequest) (*models.NotificationPreference, error) {
	if userID == 0 {
		s.log.Warn("update notification preferences unauthorized")
		return nil, dto.ErrUnauthorized
	}

	val, err := s.notificationRepo.GetNotificationPreferences(userID)
	if err != nil {
		s.log.Error(
			"failed to load notification preferences for update",
			"error", err,
			"userID", userID,
		)
		return nil, err
	}

	if req.TicketPurchased != nil {
		val.TicketPurchased = *req.TicketPurchased
	}
	if req.EventCanceled != nil {
		val.EventCanceled = *req.EventCanceled
	}
	if req.EventReminder != nil {
		val.EventReminder = *req.EventReminder
	}
	if req.PushEnabled != nil {
		val.PushEnabled = *req.PushEnabled
	}
	if req.InAppEnabled != nil {
		val.InAppEnabled = *req.InAppEnabled
	}

	if err := s.notificationRepo.UpdateNotificationPreferences(val); err != nil {
		s.log.Error(
			"failed to update notification preferences",
			"error", err,
			"userID", userID,
		)
		return nil, err
	}

	s.log.Info(
		"notification preferences updated",
		"userID", userID,
	)

	return val, nil
}

func (s *notificationService) Count(userID uint) (int64, error) {
	if userID == 0 {
		s.log.Warn("count unread notifications unauthorized")
		return 0, dto.ErrUnauthorized
	}

	count, err := s.notificationRepo.UnreadNotificationsCounts(userID)
	if err != nil {
		s.log.Error(
			"failed to count unread notifications",
			"error", err,
			"userID", userID,
		)
		return 0, err
	}

	s.log.Info(
		"unread notifications counted",
		"userID", userID,
		"count", count,
	)

	return count, nil
}
