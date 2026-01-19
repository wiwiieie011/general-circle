package services

import (
    "errors"
    e "event-service/internal/errors"
    "event-service/internal/models"
    "event-service/internal/dto"
    "io"
    "log/slog"
    "reflect"
    "testing"
)

// mockCategoryRepo is a lightweight mock for CategoryRepository used in tests.
type mockCategoryRepo struct {
    CreateFunc    func(*models.Category) error
    GetByIDFunc   func(uint) (*models.Category, error)
    DeleteFunc    func(uint) error
    GetByNameFunc func(string) (*models.Category, error)
    ListFunc      func() ([]models.Category, error)
}

func (m *mockCategoryRepo) Create(c *models.Category) error {
    if m.CreateFunc != nil {
        return m.CreateFunc(c)
    }
    return nil
}

func (m *mockCategoryRepo) GetByID(id uint) (*models.Category, error) {
    if m.GetByIDFunc != nil {
        return m.GetByIDFunc(id)
    }
    return nil, nil
}

func (m *mockCategoryRepo) Delete(id uint) error {
    if m.DeleteFunc != nil {
        return m.DeleteFunc(id)
    }
    return nil
}

func (m *mockCategoryRepo) GetByName(name string) (*models.Category, error) {
    if m.GetByNameFunc != nil {
        return m.GetByNameFunc(name)
    }
    return nil, nil
}

func (m *mockCategoryRepo) List() ([]models.Category, error) {
    if m.ListFunc != nil {
        return m.ListFunc()
    }
    return nil, nil
}

func testLogger() *slog.Logger {
    return slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelDebug}))
}

func TestCreateCategory_Success(t *testing.T) {
    repo := &mockCategoryRepo{
        GetByNameFunc: func(name string) (*models.Category, error) {
            if name != "Tech" {
                t.Fatalf("unexpected name: %s", name)
            }
            return nil, e.ErrCategoryNotFound
        },
        CreateFunc: func(c *models.Category) error {
            if c == nil || c.Name != "Tech" {
                t.Fatalf("unexpected category to create: %#v", c)
            }
            c.ID = 1
            return nil
        },
    }
    svc := NewCategoryService(repo, testLogger())

    got, err := svc.CreateCategory(dto.CreateCategoryRequest{Name: "Tech"})
    if err != nil {
        t.Fatalf("CreateCategory returned error: %v", err)
    }
    if got == nil || got.Name != "Tech" || got.ID != 1 {
        t.Fatalf("unexpected result: %#v", got)
    }
}

func TestCreateCategory_DuplicateName(t *testing.T) {
    repo := &mockCategoryRepo{
        GetByNameFunc: func(name string) (*models.Category, error) {
            return &models.Category{Base: models.Base{ID: 10}, Name: name}, nil
        },
    }
    svc := NewCategoryService(repo, testLogger())

    got, err := svc.CreateCategory(dto.CreateCategoryRequest{Name: "Tech"})
    if err == nil || !errors.Is(err, e.ErrCategoryNameExists) {
        t.Fatalf("expected ErrCategoryNameExists, got: %v", err)
    }
    if got != nil {
        t.Fatalf("expected nil category, got: %#v", got)
    }
}

func TestCreateCategory_GetByNameError(t *testing.T) {
    boom := errors.New("db error")
    repo := &mockCategoryRepo{
        GetByNameFunc: func(name string) (*models.Category, error) {
            return nil, boom
        },
    }
    svc := NewCategoryService(repo, testLogger())

    got, err := svc.CreateCategory(dto.CreateCategoryRequest{Name: "Tech"})
    if err == nil || !errors.Is(err, boom) {
        t.Fatalf("expected original error, got: %v", err)
    }
    if got != nil {
        t.Fatalf("expected nil category, got: %#v", got)
    }
}

func TestCreateCategory_CreateError(t *testing.T) {
    boom := errors.New("insert failed")
    repo := &mockCategoryRepo{
        GetByNameFunc: func(name string) (*models.Category, error) { return nil, e.ErrCategoryNotFound },
        CreateFunc:    func(c *models.Category) error { return boom },
    }
    svc := NewCategoryService(repo, testLogger())

    got, err := svc.CreateCategory(dto.CreateCategoryRequest{Name: "Tech"})
    if err == nil || !errors.Is(err, boom) {
        t.Fatalf("expected create error, got: %v", err)
    }
    if got != nil {
        t.Fatalf("expected nil category, got: %#v", got)
    }
}

func TestGetCategory_Success(t *testing.T) {
    repo := &mockCategoryRepo{
        GetByIDFunc: func(id uint) (*models.Category, error) {
            if id != 7 {
                t.Fatalf("unexpected id: %d", id)
            }
            return &models.Category{Base: models.Base{ID: id}, Name: "Tech"}, nil
        },
    }
    svc := NewCategoryService(repo, testLogger())

    got, err := svc.GetCategory(7)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if got == nil || got.ID != 7 || got.Name != "Tech" {
        t.Fatalf("unexpected result: %#v", got)
    }
}

func TestGetCategory_NotFound(t *testing.T) {
    repo := &mockCategoryRepo{
        GetByIDFunc: func(id uint) (*models.Category, error) { return nil, errors.New("not found") },
    }
    svc := NewCategoryService(repo, testLogger())

    got, err := svc.GetCategory(99)
    if err == nil || !errors.Is(err, e.ErrCategoryNotFound) {
        t.Fatalf("expected ErrCategoryNotFound, got: %v", err)
    }
    if got != nil {
        t.Fatalf("expected nil category, got: %#v", got)
    }
}

func TestDeleteCategory_Success(t *testing.T) {
    calls := struct{ deletedIDs []uint }{}
    repo := &mockCategoryRepo{
        GetByIDFunc: func(id uint) (*models.Category, error) { return &models.Category{Base: models.Base{ID: id}}, nil },
        DeleteFunc: func(id uint) error {
            calls.deletedIDs = append(calls.deletedIDs, id)
            return nil
        },
    }
    svc := NewCategoryService(repo, testLogger())

    if err := svc.DeleteCategory(5); err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if !reflect.DeepEqual(calls.deletedIDs, []uint{5}) {
        t.Fatalf("unexpected deleted ids: %#v", calls.deletedIDs)
    }
}

func TestDeleteCategory_NotFound(t *testing.T) {
    repo := &mockCategoryRepo{
        GetByIDFunc: func(id uint) (*models.Category, error) { return nil, errors.New("missing") },
    }
    svc := NewCategoryService(repo, testLogger())

    if err := svc.DeleteCategory(123); err == nil || !errors.Is(err, e.ErrCategoryNotFound) {
        t.Fatalf("expected ErrCategoryNotFound, got: %v", err)
    }
}

func TestDeleteCategory_DeleteError(t *testing.T) {
    boom := errors.New("delete failed")
    repo := &mockCategoryRepo{
        GetByIDFunc: func(id uint) (*models.Category, error) { return &models.Category{Base: models.Base{ID: id}}, nil },
        DeleteFunc:  func(id uint) error { return boom },
    }
    svc := NewCategoryService(repo, testLogger())

    if err := svc.DeleteCategory(2); err == nil || !errors.Is(err, boom) {
        t.Fatalf("expected delete error, got: %v", err)
    }
}

func TestListCategories_Success(t *testing.T) {
    expected := []models.Category{{Base: models.Base{ID: 1}, Name: "A"}, {Base: models.Base{ID: 2}, Name: "B"}}
    repo := &mockCategoryRepo{
        ListFunc: func() ([]models.Category, error) { return expected, nil },
    }
    svc := NewCategoryService(repo, testLogger())

    got, err := svc.ListCategories()
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if !reflect.DeepEqual(got, expected) {
        t.Fatalf("unexpected list result: %#v", got)
    }
}
