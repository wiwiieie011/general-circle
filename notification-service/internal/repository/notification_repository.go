package repository

import (
	"errors"
	"log/slog"
	"notification-service/internal/dto"
	"notification-service/internal/models"

	"gorm.io/gorm"
)

type NotificationRepo interface {
	Create(notification *models.Notification) error
	GetNotifications(userID uint, limit int, lastID uint) ([]models.Notification, error)
	AllRead(userID uint) error
	ReadNotificationsByID(userID, id uint) error
	DeleteNotificationsByID(userID, id uint) error
	GetNotificationPreferences(userID uint) (*models.NotificationPreference, error)
	UpdateNotificationPreferences(pref *models.NotificationPreference) error
	UnreadNotificationsCounts(userID uint) (int64, error)
}

type notificationRepo struct {
	db  *gorm.DB
	log *slog.Logger
}

func NewNotificationRepo(db *gorm.DB, log *slog.Logger) NotificationRepo {
	return &notificationRepo{
		db:  db,
		log: log,
	}
}

func (r *notificationRepo) Create(notification *models.Notification) error {

	err := r.db.Create(notification).Error

	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			r.log.Warn("notification already exists", "userID", notification.UserID, "title", notification.Title)
			return nil
		}

		r.log.Error("failed to create notification", "error", err, "userID", notification.UserID, "title", notification.Title)
		return err
	}

	r.log.Info("notification created successfully", "userID", notification.UserID, "title", notification.Title, "id", notification.ID)
	return nil
}

func (r *notificationRepo) GetNotifications(userID uint, limit int, lastID uint) ([]models.Notification, error) {
	var nots []models.Notification

	q := r.db.Select("id, title, body, read, created_at").
		Where("user_id = ?", userID).
		Order("id DESC").
		Limit(limit)

	if lastID > 0 {
		q = q.Where("id < ?", lastID)
	}

	if err := q.Find(&nots).Error; err != nil {
		r.log.Error("failed to get notifications", "error", err, "userID", userID, "limit", limit, "lastID", lastID)
		return nil, err
	}

	r.log.Info("notifications retrieved successfully", "userID", userID, "count", len(nots))
	return nots, nil
}

func (r *notificationRepo) AllRead(userID uint) error {
	if err := r.db.
		Model(&models.Notification{}).
		Where("user_id = ? AND read = ?", userID, false).
		Update("read", true).Error; err != nil {
		r.log.Error("failed to mark all notifications as read", "error", err, "userID", userID)
		return dto.ErrNotificationUpdateFailed
	}

	r.log.Info("all notifications marked as read", "userID", userID)
	return nil
}

func (r *notificationRepo) ReadNotificationsByID(userID, id uint) error {
	if err := r.db.
		Model(&models.Notification{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("read", true).Error; err != nil {
		r.log.Error("failed to mark notification as read", "error", err, "userID", userID, "notificationID", id)
		return dto.ErrNotificationUpdateFailed
	}

	r.log.Info("notification marked as read", "userID", userID, "notificationID", id)
	return nil
}

func (r *notificationRepo) DeleteNotificationsByID(userID, id uint) error {
	if err := r.db.
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&models.Notification{}).Error; err != nil {
		r.log.Error("failed to delete notification", "error", err, "userID", userID, "notificationID", id)
		return dto.ErrNotificationDeleteFailed
	}

	r.log.Info("notification deleted successfully", "userID", userID, "notificationID", id)
	return nil
}

func (r *notificationRepo) GetNotificationPreferences(userID uint) (*models.NotificationPreference, error) {
	var pref models.NotificationPreference
	if err := r.db.
		Where("user_id = ?", userID).
		First(&pref).Error; err != nil {
		r.log.Error("notification preferences not found", "error", err, "userID", userID)
		return nil, dto.ErrNotificationPreferencesNotFound
	}

	r.log.Info("notification preferences retrieved", "userID", userID)
	return &pref, nil
}

func (r *notificationRepo) UpdateNotificationPreferences(pref *models.NotificationPreference) error {
	if err := r.db.Save(pref).Error; err != nil {
		r.log.Error("failed to update notification preferences", "error", err, "userID", pref.UserID)
		return dto.ErrNotificationPreferencesUpdateFailed
	}

	r.log.Info("notification preferences updated successfully", "userID", pref.UserID)
	return nil
}

func (r *notificationRepo) UnreadNotificationsCounts(userID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&models.Notification{}).Where("user_id = ? AND read = ?", userID, false).Count(&count).Error; err != nil {
		r.log.Error("failed to count unread notifications", "error", err, "userID", userID)
		return 0, dto.ErrUnreadCountFailed
	}

	r.log.Info("unread notifications count retrieved", "userID", userID, "count", count)
	return count, nil
}
