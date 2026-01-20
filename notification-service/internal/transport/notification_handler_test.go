package transport

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"notification-service/internal/dto"
	"notification-service/internal/models"
	"notification-service/internal/services"
)

type mockService struct {
	CreateNotificationInternalFn      func(*models.Notification) error
	GetNotificationsFn                func(userID uint, limit int, lastID uint) ([]models.Notification, error)
	CheckAllFn                        func(userID uint) error
	CheckNotificationsByIDFn          func(userID, id uint) error
	DeleteNotificationByIDFn          func(userID, id uint) error
	GetNotificationPreferencesFn      func(userID uint) (*models.NotificationPreference, error)
	UpdateFn                          func(userID uint, req dto.UpdateNotificationPreferencesRequest) (*models.NotificationPreference, error)
	CountFn                           func(userID uint) (int64, error)
}

func (m *mockService) CreateNotificationInternal(n *models.Notification) error { if m.CreateNotificationInternalFn != nil { return m.CreateNotificationInternalFn(n) }; return nil }
func (m *mockService) GetNotifications(userID uint, limit int, lastID uint) ([]models.Notification, error) {
	if m.GetNotificationsFn != nil { return m.GetNotificationsFn(userID, limit, lastID) }
	return nil, nil
}
func (m *mockService) CheckAll(userID uint) error { if m.CheckAllFn != nil { return m.CheckAllFn(userID) }; return nil }
func (m *mockService) CheckNotificationsByID(userID, id uint) error {
	if m.CheckNotificationsByIDFn != nil { return m.CheckNotificationsByIDFn(userID, id) }
	return nil
}
func (m *mockService) DeleteNotificationByID(userID, id uint) error {
	if m.DeleteNotificationByIDFn != nil { return m.DeleteNotificationByIDFn(userID, id) }
	return nil
}
func (m *mockService) GetNotificationPreferences(userID uint) (*models.NotificationPreference, error) {
	if m.GetNotificationPreferencesFn != nil { return m.GetNotificationPreferencesFn(userID) }
	return &models.NotificationPreference{UserID: userID}, nil
}
func (m *mockService) Update(userID uint, req dto.UpdateNotificationPreferencesRequest) (*models.NotificationPreference, error) {
	if m.UpdateFn != nil { return m.UpdateFn(userID, req) }
	return &models.NotificationPreference{UserID: userID}, nil
}
func (m *mockService) Count(userID uint) (int64, error) { if m.CountFn != nil { return m.CountFn(userID) }; return 0, nil }

func newRouter(ms services.NotificationService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	log := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	h := NewNotificationHandler(ms, log)
	h.RegisterRoutes(r)
	return r
}

func TestGetAllNotifications(t *testing.T) {
	ms := &mockService{}
	r := newRouter(ms)

	// unauthorized
	req := httptest.NewRequest(http.MethodGet, "/notifications", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusUnauthorized, w.Code)

	// invalid limit
	req = httptest.NewRequest(http.MethodGet, "/notifications?limit=abc", nil)
	req.Header.Set("X-User-Id", "1")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)

	// invalid last_id
	req = httptest.NewRequest(http.MethodGet, "/notifications?limit=10&last_id=xyz", nil)
	req.Header.Set("X-User-Id", "1")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)

	// service error -> 500
	ms.GetNotificationsFn = func(userID uint, limit int, lastID uint) ([]models.Notification, error) { return nil, errorsNew("boom") }
	req = httptest.NewRequest(http.MethodGet, "/notifications?limit=10", nil)
	req.Header.Set("X-User-Id", "1")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusInternalServerError, w.Code)

	// success
	ms.GetNotificationsFn = func(userID uint, limit int, lastID uint) ([]models.Notification, error) {
		return []models.Notification{{Model: models.Model{ID: 2}}, {Model: models.Model{ID: 1}}}, nil
	}
	req = httptest.NewRequest(http.MethodGet, "/notifications?limit=10", nil)
	req.Header.Set("X-User-Id", "2")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestReadAllNotification(t *testing.T) {
	ms := &mockService{}
	r := newRouter(ms)

	// unauthorized
	req := httptest.NewRequest(http.MethodPut, "/notifications/read-all", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusUnauthorized, w.Code)

	// service error -> 400
	ms.CheckAllFn = func(userID uint) error { return errorsNew("fail") }
	req = httptest.NewRequest(http.MethodPut, "/notifications/read-all", nil)
	req.Header.Set("X-User-Id", "1")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)

	// success
	ms.CheckAllFn = func(userID uint) error { return nil }
	req = httptest.NewRequest(http.MethodPut, "/notifications/read-all", nil)
	req.Header.Set("X-User-Id", "1")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestReadNotificationByID(t *testing.T) {
	ms := &mockService{}
	r := newRouter(ms)

	// unauthorized
	req := httptest.NewRequest(http.MethodPut, "/notifications/12/read", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusUnauthorized, w.Code)

	// invalid id
	req = httptest.NewRequest(http.MethodPut, "/notifications/abc/read", nil)
	req.Header.Set("X-User-Id", "1")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)

	// service error
	ms.CheckNotificationsByIDFn = func(userID, id uint) error { return errorsNew("boom") }
	req = httptest.NewRequest(http.MethodPut, "/notifications/12/read", nil)
	req.Header.Set("X-User-Id", "1")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)

	// success
	ms.CheckNotificationsByIDFn = func(userID, id uint) error { return nil }
	req = httptest.NewRequest(http.MethodPut, "/notifications/12/read", nil)
	req.Header.Set("X-User-Id", "1")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteNotification(t *testing.T) {
	ms := &mockService{}
	r := newRouter(ms)

	// unauthorized
	req := httptest.NewRequest(http.MethodDelete, "/notifications/5", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusUnauthorized, w.Code)

	// invalid id
	req = httptest.NewRequest(http.MethodDelete, "/notifications/abc", nil)
	req.Header.Set("X-User-Id", "1")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)

	// service error
	ms.DeleteNotificationByIDFn = func(userID, id uint) error { return errorsNew("nope") }
	req = httptest.NewRequest(http.MethodDelete, "/notifications/5", nil)
	req.Header.Set("X-User-Id", "1")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)

	// success
	ms.DeleteNotificationByIDFn = func(userID, id uint) error { return nil }
	req = httptest.NewRequest(http.MethodDelete, "/notifications/5", nil)
	req.Header.Set("X-User-Id", "1")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestGetAndUpdatePreferences(t *testing.T) {
	ms := &mockService{}
	r := newRouter(ms)

	// get unauthorized
	req := httptest.NewRequest(http.MethodGet, "/notifications/preferences", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusUnauthorized, w.Code)

	// get service error -> 400
	ms.GetNotificationPreferencesFn = func(userID uint) (*models.NotificationPreference, error) { return nil, errorsNew("x") }
	req = httptest.NewRequest(http.MethodGet, "/notifications/preferences", nil)
	req.Header.Set("X-User-Id", "1")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)

	// get success
	ms.GetNotificationPreferencesFn = func(userID uint) (*models.NotificationPreference, error) {
		return &models.NotificationPreference{UserID: userID, PushEnabled: true}, nil
	}
	req = httptest.NewRequest(http.MethodGet, "/notifications/preferences", nil)
	req.Header.Set("X-User-Id", "1")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	// update unauthorized
	req = httptest.NewRequest(http.MethodPatch, "/notifications/preferences", bytes.NewBufferString("{}"))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusUnauthorized, w.Code)

	// invalid payload
	req = httptest.NewRequest(http.MethodPatch, "/notifications/preferences", bytes.NewBufferString("{invalid}"))
	req.Header.Set("X-User-Id", "1")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)

	// update service error -> 403
	ms.UpdateFn = func(userID uint, req dto.UpdateNotificationPreferencesRequest) (*models.NotificationPreference, error) { return nil, errorsNew("denied") }
	body, _ := json.Marshal(dto.UpdateNotificationPreferencesRequest{PushEnabled: boolPtr(true)})
	req = httptest.NewRequest(http.MethodPatch, "/notifications/preferences", bytes.NewBuffer(body))
	req.Header.Set("X-User-Id", "1")
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusForbidden, w.Code)

	// update success
	ms.UpdateFn = func(userID uint, req dto.UpdateNotificationPreferencesRequest) (*models.NotificationPreference, error) {
		return &models.NotificationPreference{UserID: userID, PushEnabled: true}, nil
	}
	body, _ = json.Marshal(dto.UpdateNotificationPreferencesRequest{PushEnabled: boolPtr(true)})
	req = httptest.NewRequest(http.MethodPatch, "/notifications/preferences", bytes.NewBuffer(body))
	req.Header.Set("X-User-Id", "1")
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestCount(t *testing.T) {
	ms := &mockService{}
	r := newRouter(ms)

	// unauthorized
	req := httptest.NewRequest(http.MethodGet, "/notifications/unread-count", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusUnauthorized, w.Code)

	// service error
	ms.CountFn = func(userID uint) (int64, error) { return 0, errorsNew("err") }
	req = httptest.NewRequest(http.MethodGet, "/notifications/unread-count", nil)
	req.Header.Set("X-User-Id", "1")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)

	// success
	ms.CountFn = func(userID uint) (int64, error) { return 3, nil }
	req = httptest.NewRequest(http.MethodGet, "/notifications/unread-count", nil)
	req.Header.Set("X-User-Id", "1")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

// helpers
func boolPtr(b bool) *bool { return &b }

// local error helper to avoid importing fmt
func errorsNew(msg string) error { return &simpleErr{msg: msg} }
type simpleErr struct{ msg string }
func (e *simpleErr) Error() string { return e.msg }
