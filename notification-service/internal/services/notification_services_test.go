package services

import (
	"errors"
	"io"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"
	"notification-service/internal/dto"
	"notification-service/internal/models"
)

type mockRepo struct {
	CreateFn                    func(*models.Notification) error
	GetNotificationsFn          func(userID uint, limit int, lastID uint) ([]models.Notification, error)
	AllReadFn                   func(userID uint) error
	ReadNotificationsByIDFn     func(userID, id uint) error
	DeleteNotificationsByIDFn   func(userID, id uint) error
	GetNotificationPreferencesFn func(userID uint) (*models.NotificationPreference, error)
	UpdateNotificationPreferencesFn func(*models.NotificationPreference) error
	UnreadNotificationsCountsFn func(userID uint) (int64, error)
}

func (m *mockRepo) Create(n *models.Notification) error { return m.callCreate(n) }
func (m *mockRepo) callCreate(n *models.Notification) error {
	if m.CreateFn != nil { return m.CreateFn(n) }
	return nil
}
func (m *mockRepo) GetNotifications(userID uint, limit int, lastID uint) ([]models.Notification, error) {
	if m.GetNotificationsFn != nil { return m.GetNotificationsFn(userID, limit, lastID) }
	return nil, nil
}
func (m *mockRepo) AllRead(userID uint) error { if m.AllReadFn != nil { return m.AllReadFn(userID) }; return nil }
func (m *mockRepo) ReadNotificationsByID(userID, id uint) error {
	if m.ReadNotificationsByIDFn != nil { return m.ReadNotificationsByIDFn(userID, id) }
	return nil
}
func (m *mockRepo) DeleteNotificationsByID(userID, id uint) error {
	if m.DeleteNotificationsByIDFn != nil { return m.DeleteNotificationsByIDFn(userID, id) }
	return nil
}
func (m *mockRepo) GetNotificationPreferences(userID uint) (*models.NotificationPreference, error) {
	if m.GetNotificationPreferencesFn != nil { return m.GetNotificationPreferencesFn(userID) }
	return &models.NotificationPreference{UserID: userID}, nil
}
func (m *mockRepo) UpdateNotificationPreferences(pref *models.NotificationPreference) error {
	if m.UpdateNotificationPreferencesFn != nil { return m.UpdateNotificationPreferencesFn(pref) }
	return nil
}
func (m *mockRepo) UnreadNotificationsCounts(userID uint) (int64, error) {
	if m.UnreadNotificationsCountsFn != nil { return m.UnreadNotificationsCountsFn(userID) }
	return 0, nil
}

func newSvc(m *mockRepo) NotificationService {
	log := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	return NewNotificationService(m, log)
}

func TestService_CreateNotificationInternal(t *testing.T) {
	m := &mockRepo{}
	svc := newSvc(m)

	// invalid user id
	err := svc.CreateNotificationInternal(&models.Notification{UserID: 0})
	require.Error(t, err)

	// success path ensures Read=false and Create called
	called := false
	m.CreateFn = func(n *models.Notification) error {
		called = true
		require.Equal(t, uint(5), n.UserID)
		require.False(t, n.Read)
		return nil
	}
	err = svc.CreateNotificationInternal(&models.Notification{UserID: 5})
	require.NoError(t, err)
	require.True(t, called)
}

func TestService_GetNotifications(t *testing.T) {
	m := &mockRepo{}
	svc := newSvc(m)

	// unauthorized
	_, err := svc.GetNotifications(0, 10, 0)
	require.ErrorIs(t, err, dto.ErrUnauthorized)

	// repo error
	m.GetNotificationsFn = func(userID uint, limit int, lastID uint) ([]models.Notification, error) {
		return nil, errors.New("db")
	}
	_, err = svc.GetNotifications(1, 10, 0)
	require.Error(t, err)

	// success
	m.GetNotificationsFn = func(userID uint, limit int, lastID uint) ([]models.Notification, error) {
		return []models.Notification{{Model: models.Model{ID: 2}}, {Model: models.Model{ID: 1}}}, nil
	}
	list, err := svc.GetNotifications(2, 10, 0)
	require.NoError(t, err)
	require.Len(t, list, 2)
}

func TestService_CheckAll(t *testing.T) {
	m := &mockRepo{}
	svc := newSvc(m)

	// unauthorized
	require.ErrorIs(t, svc.CheckAll(0), dto.ErrUnauthorized)

	// repo error
	m.AllReadFn = func(userID uint) error { return errors.New("fail") }
	require.Error(t, svc.CheckAll(1))

	// success
	m.AllReadFn = func(userID uint) error { return nil }
	require.NoError(t, svc.CheckAll(1))
}

func TestService_CheckNotificationsByID(t *testing.T) {
	m := &mockRepo{}
	svc := newSvc(m)

	require.ErrorIs(t, svc.CheckNotificationsByID(0, 1), dto.ErrUnauthorized)
	require.ErrorIs(t, svc.CheckNotificationsByID(1, 0), dto.ErrInvalidNotificationID)

	m.ReadNotificationsByIDFn = func(userID, id uint) error { return errors.New("fail") }
	require.Error(t, svc.CheckNotificationsByID(1, 2))

	m.ReadNotificationsByIDFn = func(userID, id uint) error { return nil }
	require.NoError(t, svc.CheckNotificationsByID(1, 2))
}

func TestService_DeleteNotificationByID(t *testing.T) {
	m := &mockRepo{}
	svc := newSvc(m)

	require.ErrorIs(t, svc.DeleteNotificationByID(0, 1), dto.ErrUnauthorized)
	require.ErrorIs(t, svc.DeleteNotificationByID(1, 0), dto.ErrInvalidNotificationID)

	m.DeleteNotificationsByIDFn = func(userID, id uint) error { return errors.New("fail") }
	require.Error(t, svc.DeleteNotificationByID(1, 2))

	m.DeleteNotificationsByIDFn = func(userID, id uint) error { return nil }
	require.NoError(t, svc.DeleteNotificationByID(1, 2))
}

func TestService_GetNotificationPreferences(t *testing.T) {
	m := &mockRepo{}
	svc := newSvc(m)

	_, err := svc.GetNotificationPreferences(0)
	require.ErrorIs(t, err, dto.ErrUnauthorized)

	m.GetNotificationPreferencesFn = func(userID uint) (*models.NotificationPreference, error) { return nil, errors.New("no") }
	_, err = svc.GetNotificationPreferences(1)
	require.ErrorIs(t, err, dto.ErrPreferencesNotFound)

	exp := &models.NotificationPreference{UserID: 3, PushEnabled: true}
	m.GetNotificationPreferencesFn = func(userID uint) (*models.NotificationPreference, error) { return exp, nil }
	got, err := svc.GetNotificationPreferences(3)
	require.NoError(t, err)
	require.Equal(t, exp, got)
}

func TestService_UpdatePreferences(t *testing.T) {
	m := &mockRepo{}
	svc := newSvc(m)

	_, err := svc.Update(0, dto.UpdateNotificationPreferencesRequest{})
	require.ErrorIs(t, err, dto.ErrUnauthorized)

	m.GetNotificationPreferencesFn = func(userID uint) (*models.NotificationPreference, error) { return nil, errors.New("load") }
	_, err = svc.Update(2, dto.UpdateNotificationPreferencesRequest{})
	require.Error(t, err)

	base := &models.NotificationPreference{UserID: 2, PushEnabled: true, InAppEnabled: true}
	m.GetNotificationPreferencesFn = func(userID uint) (*models.NotificationPreference, error) { return base, nil }
	updated := false
	m.UpdateNotificationPreferencesFn = func(pref *models.NotificationPreference) error {
		updated = true
		require.False(t, pref.PushEnabled)
		require.False(t, pref.InAppEnabled)
		require.True(t, pref.EventReminder)
		return nil
	}
	f := func(b bool) *bool { return &b }
	res, err := svc.Update(2, dto.UpdateNotificationPreferencesRequest{
		PushEnabled:   f(false),
		InAppEnabled:  f(false),
		EventReminder: f(true),
	})
	require.NoError(t, err)
	require.True(t, updated)
	require.False(t, res.PushEnabled)
	require.False(t, res.InAppEnabled)
	require.True(t, res.EventReminder)
}

func TestService_Count(t *testing.T) {
	m := &mockRepo{}
	svc := newSvc(m)

	_, err := svc.Count(0)
	require.ErrorIs(t, err, dto.ErrUnauthorized)

	m.UnreadNotificationsCountsFn = func(userID uint) (int64, error) { return 0, errors.New("db") }
	_, err = svc.Count(1)
	require.Error(t, err)

	m.UnreadNotificationsCountsFn = func(userID uint) (int64, error) { return 7, nil }
	got, err := svc.Count(1)
	require.NoError(t, err)
	require.Equal(t, int64(7), got)
}