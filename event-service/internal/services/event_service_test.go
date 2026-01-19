package services

import (
    "context"
    "errors"
    "event-service/internal/dto"
    e "event-service/internal/errors"
    "event-service/internal/kafka"
    "event-service/internal/models"
    "io"
    "log/slog"
    "reflect"
    "testing"
    "time"
)

type mockEventRepo struct {
    CreateFunc                func(*models.Event) error
    GetByIDFunc               func(uint) (*models.Event, error)
    UpdateFunc                func(*models.Event) error
    DeleteFunc                func(uint) error
    ListFunc                  func(dto.EventListQuery) ([]models.Event, error)
    GetByUserIDFunc           func(uint) ([]models.Event, error)
    GetEventStartingTomorrowFunc func() ([]models.Event, error)
}

func (m *mockEventRepo) Create(e *models.Event) error                          { if m.CreateFunc != nil { return m.CreateFunc(e) }; return nil }
func (m *mockEventRepo) GetByID(id uint) (*models.Event, error)                { if m.GetByIDFunc != nil { return m.GetByIDFunc(id) }; return nil, nil }
func (m *mockEventRepo) Update(e *models.Event) error                          { if m.UpdateFunc != nil { return m.UpdateFunc(e) }; return nil }
func (m *mockEventRepo) Delete(id uint) error                                  { if m.DeleteFunc != nil { return m.DeleteFunc(id) }; return nil }
func (m *mockEventRepo) List(q dto.EventListQuery) ([]models.Event, error)     { if m.ListFunc != nil { return m.ListFunc(q) }; return nil, nil }
func (m *mockEventRepo) GetByUserID(uid uint) ([]models.Event, error)          { if m.GetByUserIDFunc != nil { return m.GetByUserIDFunc(uid) }; return nil, nil }
func (m *mockEventRepo) GetEventStartingTomorrow() ([]models.Event, error)     { if m.GetEventStartingTomorrowFunc != nil { return m.GetEventStartingTomorrowFunc() }; return nil, nil }

type mockProducer struct {
    SendCancelledFunc func(context.Context, uint) error
    SendReminderFunc  func(context.Context, uint, string, time.Time) error
    CloseFunc         func() error
}

func (m *mockProducer) SendEventCancelled(ctx context.Context, eventID uint) error { if m.SendCancelledFunc != nil { return m.SendCancelledFunc(ctx, eventID) }; return nil }
func (m *mockProducer) SendEventReminder(ctx context.Context, eventID uint, title string, when time.Time) error {
    if m.SendReminderFunc != nil { return m.SendReminderFunc(ctx, eventID, title, when) }
    return nil
}
func (m *mockProducer) Close() error { if m.CloseFunc != nil { return m.CloseFunc() }; return nil }

func logger() *slog.Logger { return slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelDebug})) }

func TestEvent_Create_Success_NoCategory(t *testing.T) {
    repo := &mockEventRepo{CreateFunc: func(e *models.Event) error { e.ID = 1; return nil }}
    svc := NewEventService(repo, &mockCategoryRepo{}, &mockProducer{}, logger())

    seats := 100
    got, err := svc.CreateEvent(dto.CreateEventRequest{Title: " My Event ", UserID: 42, Seats: &seats})
    if err != nil { t.Fatalf("unexpected error: %v", err) }
    if got == nil || got.ID != 1 || got.Title != "My Event" || got.Status != string(dto.Draft) || got.UserID != 42 || got.Seats == nil || *got.Seats != 100 {
        t.Fatalf("unexpected event: %#v", got)
    }
}

func TestEvent_Create_CategoryNotFound(t *testing.T) {
    catID := uint(5)
    repo := &mockEventRepo{}
    catRepo := &mockCategoryRepo{GetByIDFunc: func(id uint) (*models.Category, error) { return nil, errors.New("missing") }}
    svc := NewEventService(repo, catRepo, &mockProducer{}, logger())

    _, err := svc.CreateEvent(dto.CreateEventRequest{Title: "Event", UserID: 1, CategoryID: &catID})
    if err == nil || !errors.Is(err, e.ErrCategoryNotFound) {
        t.Fatalf("expected ErrCategoryNotFound, got %v", err)
    }
}

func TestEvent_Create_CreateError(t *testing.T) {
    boom := errors.New("create failed")
    repo := &mockEventRepo{CreateFunc: func(e *models.Event) error { return boom }}
    svc := NewEventService(repo, &mockCategoryRepo{}, &mockProducer{}, logger())

    _, err := svc.CreateEvent(dto.CreateEventRequest{Title: "Event", UserID: 1})
    if err == nil || !errors.Is(err, boom) { t.Fatalf("expected create error, got %v", err) }
}

func TestEvent_Get_Success(t *testing.T) {
    repo := &mockEventRepo{GetByIDFunc: func(id uint) (*models.Event, error) { return &models.Event{Base: models.Base{ID: id}, Title: "E"}, nil }}
    svc := NewEventService(repo, &mockCategoryRepo{}, &mockProducer{}, logger())
    got, err := svc.GetEvent(7)
    if err != nil || got == nil || got.ID != 7 { t.Fatalf("unexpected: got=%#v err=%v", got, err) }
}

func TestEvent_Get_NotFound(t *testing.T) {
    repo := &mockEventRepo{GetByIDFunc: func(id uint) (*models.Event, error) { return nil, errors.New("missing") }}
    svc := NewEventService(repo, &mockCategoryRepo{}, &mockProducer{}, logger())
    got, err := svc.GetEvent(7)
    if err == nil || !errors.Is(err, e.ErrEventNotFound) || got != nil { t.Fatalf("expected ErrEventNotFound, got=%v", err) }
}

func TestEvent_Delete_Success(t *testing.T) {
    repo := &mockEventRepo{
        GetByIDFunc: func(id uint) (*models.Event, error) { status := string(dto.Draft); return &models.Event{Base: models.Base{ID: id}, Status: status}, nil },
        DeleteFunc:  func(id uint) error { return nil },
    }
    svc := NewEventService(repo, &mockCategoryRepo{}, &mockProducer{}, logger())
    if err := svc.DeleteEvent(3); err != nil { t.Fatalf("unexpected error: %v", err) }
}

func TestEvent_Delete_NotFound(t *testing.T) {
    repo := &mockEventRepo{GetByIDFunc: func(id uint) (*models.Event, error) { return nil, errors.New("missing") }}
    svc := NewEventService(repo, &mockCategoryRepo{}, &mockProducer{}, logger())
    if err := svc.DeleteEvent(3); err == nil || !errors.Is(err, e.ErrEventNotFound) { t.Fatalf("expected ErrEventNotFound, got %v", err) }
}

func TestEvent_Delete_NotDraft(t *testing.T) {
    repo := &mockEventRepo{GetByIDFunc: func(id uint) (*models.Event, error) { return &models.Event{Base: models.Base{ID: id}, Status: string(dto.Published)}, nil }}
    svc := NewEventService(repo, &mockCategoryRepo{}, &mockProducer{}, logger())
    if err := svc.DeleteEvent(3); err == nil || !errors.Is(err, e.ErrEventIsNotDraft) { t.Fatalf("expected ErrEventIsNotDraft, got %v", err) }
}

func TestEvent_Update_Success(t *testing.T) {
    name := " New Title "
    seats := 55
    uid := uint(9)
    repo := &mockEventRepo{
        GetByIDFunc: func(id uint) (*models.Event, error) { return &models.Event{Base: models.Base{ID: id}, Title: "old", Seats: nil, UserID: 1}, nil },
        UpdateFunc:  func(e *models.Event) error { return nil },
    }
    catID := uint(2)
    catRepo := &mockCategoryRepo{GetByIDFunc: func(id uint) (*models.Category, error) { return &models.Category{Base: models.Base{ID: id}}, nil }}
    svc := NewEventService(repo, catRepo, &mockProducer{}, logger())

    got, err := svc.UpdateEvent(dto.UpdateEventRequest{Title: &name, Seats: &seats, UserID: &uid, CategoryID: &catID}, 1)
    if err != nil { t.Fatalf("unexpected error: %v", err) }
    if got.Title != "New Title" || got.Seats == nil || *got.Seats != 55 || got.UserID != 9 || got.CategoryID == nil || *got.CategoryID != 2 { t.Fatalf("unexpected event: %#v", got) }
}

func TestEvent_Update_NotFound(t *testing.T) {
    repo := &mockEventRepo{GetByIDFunc: func(id uint) (*models.Event, error) { return nil, errors.New("missing") }}
    svc := NewEventService(repo, &mockCategoryRepo{}, &mockProducer{}, logger())
    _, err := svc.UpdateEvent(dto.UpdateEventRequest{}, 1)
    if err == nil || !errors.Is(err, e.ErrEventNotFound) { t.Fatalf("expected ErrEventNotFound, got %v", err) }
}

func TestEvent_Update_EmptyTitle(t *testing.T) {
    repo := &mockEventRepo{GetByIDFunc: func(id uint) (*models.Event, error) { return &models.Event{Base: models.Base{ID: id}, Title: "t"}, nil }}
    svc := NewEventService(repo, &mockCategoryRepo{}, &mockProducer{}, logger())
    empty := "  "
    _, err := svc.UpdateEvent(dto.UpdateEventRequest{Title: &empty}, 1)
    if err == nil || !errors.Is(err, e.ErrEmptyTitle) { t.Fatalf("expected ErrEmptyTitle, got %v", err) }
}

func TestEvent_Update_BadSeats(t *testing.T) {
    repo := &mockEventRepo{GetByIDFunc: func(id uint) (*models.Event, error) { return &models.Event{Base: models.Base{ID: id}}, nil }}
    svc := NewEventService(repo, &mockCategoryRepo{}, &mockProducer{}, logger())
    seats := -1
    _, err := svc.UpdateEvent(dto.UpdateEventRequest{Seats: &seats}, 1)
    if err == nil || !errors.Is(err, e.ErrNotCorrectNum) { t.Fatalf("expected ErrNotCorrectNum, got %v", err) }
}

func TestEvent_Update_BadCategory(t *testing.T) {
    repo := &mockEventRepo{GetByIDFunc: func(id uint) (*models.Event, error) { return &models.Event{Base: models.Base{ID: id}}, nil }}
    catRepo := &mockCategoryRepo{GetByIDFunc: func(id uint) (*models.Category, error) { return nil, errors.New("missing") }}
    svc := NewEventService(repo, catRepo, &mockProducer{}, logger())
    catID := uint(77)
    _, err := svc.UpdateEvent(dto.UpdateEventRequest{CategoryID: &catID}, 1)
    if err == nil || !errors.Is(err, e.ErrCategoryNotFound) { t.Fatalf("expected ErrCategoryNotFound, got %v", err) }
}

func TestEvent_List_Success(t *testing.T) {
    want := []models.Event{{Base: models.Base{ID: 1}}, {Base: models.Base{ID: 2}}}
    repo := &mockEventRepo{ListFunc: func(q dto.EventListQuery) ([]models.Event, error) { return want, nil }}
    svc := NewEventService(repo, &mockCategoryRepo{}, &mockProducer{}, logger())
    got, err := svc.ListEvents(dto.EventListQuery{})
    if err != nil || !reflect.DeepEqual(got, want) { t.Fatalf("unexpected: %#v err=%v", got, err) }
}

func TestEvent_Publish_Success(t *testing.T) {
    updated := false
    repo := &mockEventRepo{
        GetByIDFunc: func(id uint) (*models.Event, error) { return &models.Event{Base: models.Base{ID: id}, Status: string(dto.Draft)}, nil },
        UpdateFunc:  func(e *models.Event) error { updated = true; if e.Status != string(dto.Published) { t.Fatalf("status not updated: %s", e.Status) }; return nil },
    }
    svc := NewEventService(repo, &mockCategoryRepo{}, &mockProducer{}, logger())
    if err := svc.PublishEvent(1); err != nil { t.Fatalf("unexpected error: %v", err) }
    if !updated { t.Fatalf("expected update to be called") }
}

func TestEvent_Publish_NotFound(t *testing.T) {
    repo := &mockEventRepo{GetByIDFunc: func(id uint) (*models.Event, error) { return nil, errors.New("missing") }}
    svc := NewEventService(repo, &mockCategoryRepo{}, &mockProducer{}, logger())
    if err := svc.PublishEvent(1); err == nil || !errors.Is(err, e.ErrEventNotFound) { t.Fatalf("expected ErrEventNotFound, got %v", err) }
}

func TestEvent_Publish_NotDraft(t *testing.T) {
    repo := &mockEventRepo{GetByIDFunc: func(id uint) (*models.Event, error) { return &models.Event{Base: models.Base{ID: id}, Status: string(dto.Published)}, nil }}
    svc := NewEventService(repo, &mockCategoryRepo{}, &mockProducer{}, logger())
    if err := svc.PublishEvent(1); err == nil || !errors.Is(err, e.ErrEventIsNotDraft) { t.Fatalf("expected ErrEventIsNotDraft, got %v", err) }
}

func TestEvent_Cancel_Success_ProducerErrorIgnored(t *testing.T) {
    updated := false
    repo := &mockEventRepo{
        GetByIDFunc: func(id uint) (*models.Event, error) { return &models.Event{Base: models.Base{ID: id}, Status: string(dto.Published)}, nil },
        UpdateFunc:  func(e *models.Event) error { updated = true; if e.Status != string(dto.Cancelled) { t.Fatalf("status not set to cancelled") }; return nil },
    }
    prod := &mockProducer{SendCancelledFunc: func(ctx context.Context, id uint) error { return errors.New("kafka down") }}
    svc := NewEventService(repo, &mockCategoryRepo{}, prod, logger())
    if err := svc.CancelEvent(1); err != nil { t.Fatalf("unexpected error: %v", err) }
    if !updated { t.Fatalf("expected update to be called") }
}

func TestEvent_Cancel_NotFound(t *testing.T) {
    repo := &mockEventRepo{GetByIDFunc: func(id uint) (*models.Event, error) { return nil, errors.New("missing") }}
    svc := NewEventService(repo, &mockCategoryRepo{}, &mockProducer{}, logger())
    if err := svc.CancelEvent(1); err == nil || !errors.Is(err, e.ErrEventNotFound) { t.Fatalf("expected ErrEventNotFound, got %v", err) }
}

func TestEvent_Cancel_NotPublished(t *testing.T) {
    repo := &mockEventRepo{GetByIDFunc: func(id uint) (*models.Event, error) { return &models.Event{Base: models.Base{ID: id}, Status: string(dto.Draft)}, nil }}
    svc := NewEventService(repo, &mockCategoryRepo{}, &mockProducer{}, logger())
    if err := svc.CancelEvent(1); err == nil || !errors.Is(err, e.ErrEventIsNotPublished) { t.Fatalf("expected ErrEventIsNotPublished, got %v", err) }
}

func TestEvent_GetByUserID_Success(t *testing.T) {
    want := []models.Event{{Base: models.Base{ID: 1}}}
    repo := &mockEventRepo{GetByUserIDFunc: func(uid uint) ([]models.Event, error) { return want, nil }}
    svc := NewEventService(repo, &mockCategoryRepo{}, &mockProducer{}, logger())
    got, err := svc.GetEventsByUserID(42)
    if err != nil || !reflect.DeepEqual(got, want) { t.Fatalf("unexpected: %#v err=%v", got, err) }
}

func TestEvent_SendReminders(t *testing.T) {
    now := time.Now()
    events := []models.Event{
        {Base: models.Base{ID: 1}, Title: "No schedule", Schedule: nil},
        {Base: models.Base{ID: 2}, Title: "Two items", Schedule: []models.EventSchedule{{StartAt: now.Add(2*time.Hour)}, {StartAt: now.Add(1 * time.Hour)}}},
    }
    called := struct{ ids []uint }{}
    repo := &mockEventRepo{GetEventStartingTomorrowFunc: func() ([]models.Event, error) { return events, nil }}
    prod := &mockProducer{SendReminderFunc: func(ctx context.Context, id uint, title string, when time.Time) error { called.ids = append(called.ids, id); return nil }}
    svc := NewEventService(repo, &mockCategoryRepo{}, prod, logger())
    if err := svc.SendEventReminders(context.Background()); err != nil { t.Fatalf("unexpected error: %v", err) }
    if !reflect.DeepEqual(called.ids, []uint{2}) { t.Fatalf("expected reminder for event 2, got %#v", called.ids) }
}

func TestEvent_SendReminders_RepoError(t *testing.T) {
    repo := &mockEventRepo{GetEventStartingTomorrowFunc: func() ([]models.Event, error) { return nil, errors.New("db") }}
    svc := NewEventService(repo, &mockCategoryRepo{}, &mockProducer{}, logger())
    if err := svc.SendEventReminders(context.Background()); err == nil { t.Fatalf("expected error") }
}

// Ensure mockProducer satisfies interface
var _ kafka.EventProducer = (*mockProducer)(nil)
