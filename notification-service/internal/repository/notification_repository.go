package repository

import (
	"log/slog"
	"notification-service/internal/dto"
	"notification-service/internal/models"

	"gorm.io/gorm"
)

type NotificationRepo interface {
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
		return nil, dto.ErrNotificationsNotFound
	}

	return nots, nil
}

func (r *notificationRepo) AllRead(userID uint) error {
	var list []models.Notification
	if err := r.db.
		Model(&list).
		Where("user_id = ? AND read = ?", userID, false).
		Update("read", true).Error; err != nil {
		return dto.ErrNotificationUpdateFailed
	}

	return nil
}

func (r *notificationRepo) ReadNotificationsByID(userID, id uint) error {
	var not models.Notification
	if err := r.db.
		Model(&not).
		Where("id = ? AND user_id = ?", id, userID).
		Update("read", true).Error; err != nil {
		return dto.ErrNotificationUpdateFailed
	}

	return nil
}

func (r *notificationRepo) DeleteNotificationsByID(userID, id uint) error {
	var not models.Notification
	if err := r.db.
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&not).Error; err != nil {
		return dto.ErrNotificationDeleteFailed
	}

	return nil
}

func (r *notificationRepo) GetNotificationPreferences(userID uint) (*models.NotificationPreference, error) {
	var pref models.NotificationPreference
	if err := r.db.
		Where("user_id = ?", userID).
		First(&pref).Error; err != nil {
		return nil, dto.ErrNotificationPreferencesNotFound
	}
	return &pref, nil
}

func (r *notificationRepo) UpdateNotificationPreferences(pref *models.NotificationPreference) error {
	if err := r.db.Save(pref).Error; err != nil {
		return dto.ErrNotificationPreferencesUpdateFailed
	}
	return nil
}

func (r *notificationRepo) UnreadNotificationsCounts(userID uint) (int64, error) {
	var list []models.Notification
	var count int64
	if err := r.db.Model(&list).Where("user_id = ? AND read = ?", userID, false).Count(&count).Error; err != nil {
		return count, err
	}
	return count, dto.ErrUnreadCountFailed
}
