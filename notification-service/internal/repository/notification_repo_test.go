package repository

import (
	"io"
	"log/slog"
	"testing"

	"notification-service/internal/models"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func newTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
}

func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&models.Notification{}, &models.NotificationPreference{}))
	return db
}

func TestNotificationRepo_Create(t *testing.T) {
	db := newTestDB(t)
	repo := NewNotificationRepo(db, newTestLogger())

	n := &models.Notification{UserID: 1, Title: "A", Body: "a body", Read: false}
	require.NoError(t, repo.Create(n))
	require.NotZero(t, n.ID)

	// verify persisted
	var got models.Notification
	err := db.First(&got, n.ID).Error
	require.NoError(t, err)
	require.Equal(t, "A", got.Title)
}

func TestNotificationRepo_GetNotifications(t *testing.T) {
	db := newTestDB(t)
	repo := NewNotificationRepo(db, newTestLogger())


	n1 := &models.Notification{UserID: 1, Title: "A"}
	n2 := &models.Notification{UserID: 1, Title: "B"}
	n3 := &models.Notification{UserID: 2, Title: "C"}
	require.NoError(t, repo.Create(n1))
	require.NoError(t, repo.Create(n2))
	require.NoError(t, repo.Create(n3))

	list, err := repo.GetNotifications(1, 10, 0)
	require.NoError(t, err)
	require.Len(t, list, 2)
	require.Greater(t, list[0].ID, list[1].ID)

	lastID := list[0].ID
	list2, err := repo.GetNotifications(1, 10, lastID)
	require.NoError(t, err)
	require.Len(t, list2, 1)
	require.Equal(t, list[1].ID, list2[0].ID)
}

func TestNotificationRepo_ReadNotificationsByID(t *testing.T) {
	db := newTestDB(t)
	repo := NewNotificationRepo(db, newTestLogger())

	n1 := &models.Notification{UserID: 1, Title: "A", Read: false}
	n2 := &models.Notification{UserID: 1, Title: "B", Read: false}
	require.NoError(t, repo.Create(n1))
	require.NoError(t, repo.Create(n2))

	
	cnt, err := repo.UnreadNotificationsCounts(1)
	require.NoError(t, err)
	require.Equal(t, int64(2), cnt)

	
	require.NoError(t, repo.ReadNotificationsByID(1, n1.ID))
	cnt, err = repo.UnreadNotificationsCounts(1)
	require.NoError(t, err)
	require.Equal(t, int64(1), cnt)
}

func TestNotificationRepo_AllRead(t *testing.T) {
	db := newTestDB(t)
	repo := NewNotificationRepo(db, newTestLogger())

	n1 := &models.Notification{UserID: 1, Title: "A", Read: false}
	n2 := &models.Notification{UserID: 1, Title: "B", Read: false}
	require.NoError(t, repo.Create(n1))
	require.NoError(t, repo.Create(n2))

	require.NoError(t, repo.AllRead(1))
	cnt, err := repo.UnreadNotificationsCounts(1)
	require.NoError(t, err)
	require.Equal(t, int64(0), cnt)
}

func TestNotificationRepo_DeleteNotificationsByID(t *testing.T) {
	db := newTestDB(t)
	repo := NewNotificationRepo(db, newTestLogger())

	n1 := &models.Notification{UserID: 1, Title: "A"}
	n2 := &models.Notification{UserID: 1, Title: "B"}
	require.NoError(t, repo.Create(n1))
	require.NoError(t, repo.Create(n2))

	require.NoError(t, repo.DeleteNotificationsByID(1, n2.ID))
	list, err := repo.GetNotifications(1, 10, 0)
	require.NoError(t, err)
	require.Len(t, list, 1)
	require.Equal(t, n1.ID, list[0].ID)
}

func TestNotificationRepo_Preferences_CreateDefault(t *testing.T) {
	db := newTestDB(t)
	repo := NewNotificationRepo(db, newTestLogger())

	
	pref, err := repo.GetNotificationPreferences(10)
	require.NoError(t, err)
	require.NotNil(t, pref)
	require.Equal(t, uint(10), pref.UserID)
	require.True(t, pref.TicketPurchased)
	require.True(t, pref.EventCanceled)
	require.True(t, pref.EventReminder)
	require.True(t, pref.PushEnabled)
	require.True(t, pref.InAppEnabled)
}

func TestNotificationRepo_Preferences_Update(t *testing.T) {
	db := newTestDB(t)
	repo := NewNotificationRepo(db, newTestLogger())

	pref, err := repo.GetNotificationPreferences(10)
	require.NoError(t, err)

	pref.PushEnabled = false
	pref.InAppEnabled = false
	pref.TicketPurchased = false
	require.NoError(t, repo.UpdateNotificationPreferences(pref))

	got, err := repo.GetNotificationPreferences(10)
	require.NoError(t, err)
	require.False(t, got.PushEnabled)
	require.False(t, got.InAppEnabled)
	require.False(t, got.TicketPurchased)
}
