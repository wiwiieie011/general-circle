package services

import (
	"errors"
	"log/slog"
	"notification-service/internal/dto"
	"notification-service/internal/models"
	"notification-service/internal/repository"
)

type NotificationService interface {
	GetNotifications(userID uint, limit int, lastID uint) ([]models.Notification, error)
	CheckAll(userID uint) error
	CheckNotificationsByID(userID, id uint) error
	DeleteNotificationByID(userID, id uint) error
	GetNotificationPreferences(userID uint) (*models.NotificationPreference, error)
	Update(userID uint, req dto.UpdateNotificationPreferencesRequest) (*models.NotificationPreference, error)
	Count(userID uint) (int64, error)
}

type notificationService struct {
	notificationRepo repository.NotificationRepo
	log              *slog.Logger
}

func NewNotifictaonService(notificationRepo repository.NotificationRepo, log *slog.Logger) NotificationService {
	return &notificationService{
		notificationRepo: notificationRepo,
		log:              log,
	}
}

func (s *notificationService) GetNotifications(userID uint, limit int, lastID uint) ([]models.Notification, error) {
	nots, err := s.notificationRepo.GetNotifications(userID, limit, lastID)

	if err != nil {
		return []models.Notification{}, err
	}

	return nots, nil
}

func (s *notificationService) CheckAll(userID uint) error {
	if userID == 0 {
		return dto.ErrUnauthorized
	}

	return s.notificationRepo.AllRead(userID)
}

func (s *notificationService) CheckNotificationsByID(userID, id uint) error {
	if userID == 0 {
		return dto.ErrUnauthorized
	}

	if id == 0 {
		return dto.ErrInvalidNotificationID
	}

	return s.notificationRepo.ReadNotificationsByID(userID, id)
}

func (s *notificationService) DeleteNotificationByID(userID, id uint) error {
	if userID == 0 {
		return dto.ErrUnauthorized
	}

	if id == 0 {
		return dto.ErrInvalidNotificationID
	}

	return s.notificationRepo.DeleteNotificationsByID(userID, id)
}

func (s *notificationService) GetNotificationPreferences(userID uint) (*models.NotificationPreference, error) {
	if userID == 0 {
		return nil, dto.ErrUnauthorized
	}

	val, err := s.notificationRepo.GetNotificationPreferences(userID)
	if err != nil {
		return nil, dto.ErrPreferencesNotFound
	}

	return val, nil
}

func (s *notificationService) Update(userID uint, req dto.UpdateNotificationPreferencesRequest) (*models.NotificationPreference, error) {

	if userID == 0 {
		return nil, dto.ErrUnauthorized
	}

	val, err := s.notificationRepo.GetNotificationPreferences(userID)

	if err != nil {
		return nil, errors.New("err")
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
		return nil, err
	}

	return val, nil
}

func (s *notificationService) Count(userID uint) (int64, error) {
	if userID == 0 {
		return 0, dto.ErrUnauthorized
	}

	count, err := s.notificationRepo.UnreadNotificationsCounts(userID)
	if err != nil {
		return 0, err
	}

	return count, nil
}
