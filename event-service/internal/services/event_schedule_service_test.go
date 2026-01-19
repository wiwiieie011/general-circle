package services

import (
	"errors"
	"event-service/internal/dto"
	e "event-service/internal/errors"
	"event-service/internal/models"
	"testing"
	"time"
)

type mockEventScheduleRepo struct {
	CreateFunc       func(*models.EventSchedule) error
	GetByIDFunc      func(uint) (*models.EventSchedule, error)
	GetByEventIDFunc func(uint) ([]models.EventSchedule, error)
}

func (m *mockEventScheduleRepo) Create(s *models.EventSchedule) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(s)
	}
	return nil
}

func (m *mockEventScheduleRepo) GetByID(id uint) (*models.EventSchedule, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(id)
	}
	return nil, nil
}

func (m *mockEventScheduleRepo) GetByEventID(eid uint) ([]models.EventSchedule, error) {
	if m.GetByEventIDFunc != nil {
		return m.GetByEventIDFunc(eid)
	}
	return nil, nil
}

func TestSchedule_GetByEventID_Success(t *testing.T) {
	repo := &mockEventScheduleRepo{GetByEventIDFunc: func(eid uint) ([]models.EventSchedule, error) { return []models.EventSchedule{{}, {}}, nil }}
	evtRepo := &mockEventRepo{GetByIDFunc: func(id uint) (*models.Event, error) { return &models.Event{Base: models.Base{ID: id}}, nil }}

	svc := NewEventScheduleService(repo, evtRepo, logger())

	got, err := svc.GetScheduleByEventID(1)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 schedules, got %d", len(got))
	}
}

func TestSchedule_GetByEventID_EventNotFound(t *testing.T) {
	repo := &mockEventScheduleRepo{}
	evtRepo := &mockEventRepo{GetByIDFunc: func(id uint) (*models.Event, error) { return nil, errors.New("missing") }}
	svc := NewEventScheduleService(repo, evtRepo, logger())
	_, err := svc.GetScheduleByEventID(1)
	if err == nil || !errors.Is(err, e.ErrEventNotFound) {
		t.Fatalf("expected ErrEventNotFound, got %v", err)
	}
}

func TestSchedule_Create_Success(t *testing.T) {
	now := time.Now()
	repo := &mockEventScheduleRepo{CreateFunc: func(s *models.EventSchedule) error {
		s.ID = 10
		return nil
	}}
	evtRepo := &mockEventRepo{GetByIDFunc: func(id uint) (*models.Event, error) { return &models.Event{Base: models.Base{ID: id}}, nil }}

	svc := NewEventScheduleService(repo, evtRepo, logger())

	got, err := svc.CreateScheduleForEvent(2, dto.CreateScheduleRequest{ActivityName: "Talk", Speaker: "Alice", StartAt: now, EndAt: now.Add(time.Hour)})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil {
		t.Fatalf("expected non-nil schedule")
	}
	if got.ID != 10 {
		t.Fatalf("expected ID=10, got %d", got.ID)
	}
	if got.EventID != 2 {
		t.Fatalf("expected EventID=2, got %d", got.EventID)
	}
	if got.ActivityName != "Talk" {
		t.Fatalf("expected ActivityName 'Talk', got %q", got.ActivityName)
	}
	if got.Speaker != "Alice" {
		t.Fatalf("expected Speaker 'Alice', got %q", got.Speaker)
	}
}

func TestSchedule_Create_EventNotFound(t *testing.T) {
	now := time.Now()
	repo := &mockEventScheduleRepo{}
	evtRepo := &mockEventRepo{GetByIDFunc: func(id uint) (*models.Event, error) { return nil, errors.New("missing") }}
	svc := NewEventScheduleService(repo, evtRepo, logger())
	_, err := svc.CreateScheduleForEvent(2, dto.CreateScheduleRequest{ActivityName: "Talk", Speaker: "Alice", StartAt: now, EndAt: now.Add(time.Hour)})
	if err == nil || !errors.Is(err, e.ErrEventNotFound) {
		t.Fatalf("expected ErrEventNotFound, got %v", err)
	}
}

func TestSchedule_Create_BadTime(t *testing.T) {
	now := time.Now()
	repo := &mockEventScheduleRepo{}
	evtRepo := &mockEventRepo{GetByIDFunc: func(id uint) (*models.Event, error) { return &models.Event{Base: models.Base{ID: id}}, nil }}
	svc := NewEventScheduleService(repo, evtRepo, logger())
	_, err := svc.CreateScheduleForEvent(2, dto.CreateScheduleRequest{ActivityName: "Talk", Speaker: "Alice", StartAt: now, EndAt: now.Add(-time.Hour)})
	if err == nil || !errors.Is(err, e.ErrNotCorrectScheduleTime) {
		t.Fatalf("expected ErrNotCorrectScheduleTime, got %v", err)
	}
}
