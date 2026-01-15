package dto

import "errors"

//repo
var (
    ErrNotificationsNotFound       = errors.New("уведомления не найдены")
    ErrNotificationNotFound        = errors.New("уведомление не найдено")
    ErrNotificationUpdateFailed    = errors.New("не удалось обновить уведомления")
    ErrNotificationDeleteFailed    = errors.New("не удалось удалить уведомление")
    ErrNotificationPreferencesNotFound = errors.New("настройки уведомлений не найдены")
    ErrNotificationPreferencesUpdateFailed = errors.New("не удалось обновить настройки уведомлений")
    ErrUnreadCountFailed           = errors.New("не удалось посчитать непрочитанные уведомления")
)

var (
	ErrUnauthorized          = errors.New("unauthorized")
	ErrInvalidNotificationID = errors.New("invalid notification id")
	ErrPreferencesNotFound   = errors.New("notification preferences not found")
)