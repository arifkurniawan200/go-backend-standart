package repository

import (
	"context"
	"testing"

	"github.com/arifkurniawan200/go-backend-standart/internal/domain"
)

func TestCreate(t *testing.T) {
	repo := NewUserRepository()
	ctx := context.Background()

	user := &domain.User{
		Name:     "Alice",
		Email:    "alice@example.com",
		Password: "hashed123",
	}

	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user.ID == "" {
		t.Error("expected user ID to be set")
	}
	if user.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
}

func TestCreate_Duplicate(t *testing.T) {
	repo := NewUserRepository()
	ctx := context.Background()

	user := &domain.User{ID: "fixed-id", Name: "Alice", Email: "alice@example.com"}

	err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("first create should succeed: %v", err)
	}

	user2 := &domain.User{ID: user.ID, Name: "Bob", Email: "bob@example.com"}
	err = repo.Create(ctx, user2)
	if err == nil {
		t.Error("expected ErrUserExists on duplicate ID")
	}
}

func TestFindByID(t *testing.T) {
	repo := NewUserRepository()
	ctx := context.Background()

	user := &domain.User{Name: "Alice", Email: "alice@example.com"}
	repo.Create(ctx, user)

	found, err := repo.FindByID(ctx, user.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if found.Name != "Alice" {
		t.Errorf("expected Name Alice, got %s", found.Name)
	}
}

func TestFindByID_NotFound(t *testing.T) {
	repo := NewUserRepository()
	ctx := context.Background()

	_, err := repo.FindByID(ctx, "nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent user")
	}
	if err != ErrUserNotFound {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}

func TestFindByEmail(t *testing.T) {
	repo := NewUserRepository()
	ctx := context.Background()

	repo.Create(ctx, &domain.User{Name: "Alice", Email: "alice@example.com"})
	repo.Create(ctx, &domain.User{Name: "Bob", Email: "bob@example.com"})

	found, err := repo.FindByEmail(ctx, "bob@example.com")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if found.Name != "Bob" {
		t.Errorf("expected Name Bob, got %s", found.Name)
	}
}

func TestFindByEmail_NotFound(t *testing.T) {
	repo := NewUserRepository()
	ctx := context.Background()

	_, err := repo.FindByEmail(ctx, "nonexistent@example.com")
	if err != ErrUserNotFound {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}

func TestUpdate(t *testing.T) {
	repo := NewUserRepository()
	ctx := context.Background()

	user := &domain.User{Name: "Alice", Email: "alice@example.com"}
	repo.Create(ctx, user)

	user.Name = "Alice Updated"
	err := repo.Update(ctx, user)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	found, _ := repo.FindByID(ctx, user.ID)
	if found.Name != "Alice Updated" {
		t.Errorf("expected Name Alice Updated, got %s", found.Name)
	}
}

func TestUpdate_NotFound(t *testing.T) {
	repo := NewUserRepository()
	ctx := context.Background()

	user := &domain.User{ID: "nonexistent", Name: "Ghost"}
	err := repo.Update(ctx, user)
	if err != ErrUserNotFound {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}

func TestDelete(t *testing.T) {
	repo := NewUserRepository()
	ctx := context.Background()

	user := &domain.User{Name: "Alice", Email: "alice@example.com"}
	repo.Create(ctx, user)

	err := repo.Delete(ctx, user.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = repo.FindByID(ctx, user.ID)
	if err != ErrUserNotFound {
		t.Errorf("expected ErrUserNotFound after delete, got %v", err)
	}
}

func TestDelete_NotFound(t *testing.T) {
	repo := NewUserRepository()
	ctx := context.Background()

	err := repo.Delete(ctx, "nonexistent")
	if err != ErrUserNotFound {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}

func TestList(t *testing.T) {
	repo := NewUserRepository()
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		repo.Create(ctx, &domain.User{Name: "User", Email: "user@example.com"})
	}

	users, err := repo.List(ctx, 3, 0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(users) != 3 {
		t.Errorf("expected 3 users with limit=3, got %d", len(users))
	}

	users2, err := repo.List(ctx, 10, 3)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(users2) != 2 {
		t.Errorf("expected 2 users with offset=3 (5 total), got %d", len(users2))
	}
}
