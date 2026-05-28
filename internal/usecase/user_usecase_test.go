package usecase

import (
	"context"
	"testing"

	"github.com/arifkurniawan200/go-backend-standart/internal/domain"
	"github.com/arifkurniawan200/go-backend-standart/internal/repository"
)

type mockRepo struct {
	users map[string]*domain.User
}

func newMockRepo() *mockRepo {
	return &mockRepo{users: make(map[string]*domain.User)}
}

func (m *mockRepo) Create(ctx context.Context, user *domain.User) error {
	user.ID = "generated-id"
	m.users[user.ID] = user
	return nil
}
func (m *mockRepo) FindByID(ctx context.Context, id string) (*domain.User, error) {
	if u, ok := m.users[id]; ok {
		return u, nil
	}
	return nil, repository.ErrUserNotFound
}
func (m *mockRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	for _, u := range m.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, repository.ErrUserNotFound
}
func (m *mockRepo) Update(ctx context.Context, user *domain.User) error {
	if _, ok := m.users[user.ID]; !ok {
		return repository.ErrUserNotFound
	}
	m.users[user.ID] = user
	return nil
}
func (m *mockRepo) Delete(ctx context.Context, id string) error {
	if _, ok := m.users[id]; !ok {
		return repository.ErrUserNotFound
	}
	delete(m.users, id)
	return nil
}
func (m *mockRepo) List(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	var result []*domain.User
	count := 0
	for _, u := range m.users {
		if count >= offset && len(result) < limit {
			result = append(result, u)
		}
		count++
	}
	return result, nil
}

func TestUsecase_Create(t *testing.T) {
	uc := NewUserUsecase(newMockRepo())
	ctx := context.Background()

	req := &domain.UserCreateRequest{
		Name:     "Alice",
		Email:    "alice@example.com",
		Password: "password123",
	}

	user, err := uc.Create(ctx, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user.ID == "" {
		t.Error("expected user ID to be set")
	}
	if user.Name != "Alice" {
		t.Errorf("expected Name Alice, got %s", user.Name)
	}
}

func TestUsecase_Create_InvalidInput(t *testing.T) {
	uc := NewUserUsecase(newMockRepo())
	ctx := context.Background()

	_, err := uc.Create(ctx, &domain.UserCreateRequest{})
	if err != ErrInvalidInput {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}
}

func TestUsecase_GetByID(t *testing.T) {
	repo := newMockRepo()
	uc := NewUserUsecase(repo)
	ctx := context.Background()

	repo.Create(ctx, &domain.User{Name: "Alice", Email: "alice@example.com"})

	user, err := uc.GetByID(ctx, "generated-id")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user.Name != "Alice" {
		t.Errorf("expected Name Alice, got %s", user.Name)
	}
}

func TestUsecase_GetByID_NotFound(t *testing.T) {
	uc := NewUserUsecase(newMockRepo())
	ctx := context.Background()

	_, err := uc.GetByID(ctx, "nonexistent")
	if err != repository.ErrUserNotFound {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}

func TestUsecase_Update(t *testing.T) {
	repo := newMockRepo()
	uc := NewUserUsecase(repo)
	ctx := context.Background()

	repo.Create(ctx, &domain.User{Name: "Alice", Email: "alice@example.com"})

	req := &domain.UserUpdateRequest{
		ID:   "generated-id",
		Name: "Alice Updated",
	}

	user, err := uc.Update(ctx, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user.Name != "Alice Updated" {
		t.Errorf("expected Name Alice Updated, got %s", user.Name)
	}
}

func TestUsecase_Delete(t *testing.T) {
	repo := newMockRepo()
	uc := NewUserUsecase(repo)
	ctx := context.Background()

	repo.Create(ctx, &domain.User{Name: "Alice", Email: "alice@example.com"})

	err := uc.Delete(ctx, "generated-id")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = uc.GetByID(ctx, "generated-id")
	if err != repository.ErrUserNotFound {
		t.Errorf("expected ErrUserNotFound after delete, got %v", err)
	}
}

func TestUsecase_List(t *testing.T) {
	repo := newMockRepo()
	uc := NewUserUsecase(repo)
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		repo.Create(ctx, &domain.User{Name: "User", Email: "user@example.com"})
	}

	users, err := uc.List(ctx, 3, 0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(users) != 3 {
		t.Errorf("expected 3 users, got %d", len(users))
	}
}

func TestUsecase_List_ZeroLimit(t *testing.T) {
	repo := newMockRepo()
	uc := NewUserUsecase(repo)
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		repo.Create(ctx, &domain.User{Name: "User", Email: "user@example.com"})
	}

	users, err := uc.List(ctx, 0, 0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	// Should default to 10 when limit <= 0
	if len(users) != 5 {
		t.Errorf("expected 5 users (all with default limit), got %d", len(users))
	}
}
