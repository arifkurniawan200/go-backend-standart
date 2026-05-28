package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arifkurniawan200/go-backend-standart/internal/domain"
	"github.com/arifkurniawan200/go-backend-standart/internal/usecase"
	"go.uber.org/zap/zaptest"
)

type mockUserRepo struct {
	users map[string]*domain.User
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{users: make(map[string]*domain.User)}
}

func (m *mockUserRepo) Create(ctx context.Context, user *domain.User) error {
	user.ID = "mock-id"
	m.users[user.ID] = user
	return nil
}
func (m *mockUserRepo) FindByID(ctx context.Context, id string) (*domain.User, error) { return nil, nil }
func (m *mockUserRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) { return nil, nil }
func (m *mockUserRepo) Update(ctx context.Context, user *domain.User) error           { return nil }
func (m *mockUserRepo) Delete(ctx context.Context, id string) error                    { return nil }
func (m *mockUserRepo) List(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	return nil, nil
}

func setupHandler() *UserHandler {
	repo := newMockUserRepo()
	uc := usecase.NewUserUsecase(repo)
	log := zaptest.NewLogger(nil)
	return NewUserHandler(uc, log)
}

func TestCreate_ValidationErrors(t *testing.T) {
	tests := []struct {
		name    string
		body    any
		wantErr string
	}{
		{
			name: "missing name",
			body: domain.UserCreateRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			wantErr: "name",
		},
		{
			name: "empty email",
			body: domain.UserCreateRequest{
				Name:     "John Doe",
				Email:    "",
				Password: "password123",
			},
			wantErr: "email",
		},
		{
			name: "invalid email format",
			body: domain.UserCreateRequest{
				Name:     "John Doe",
				Email:    "not-an-email",
				Password: "password123",
			},
			wantErr: "email",
		},
		{
			name: "short password",
			body: domain.UserCreateRequest{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "short",
			},
			wantErr: "password",
		},
		{
			name: "name too short",
			body: domain.UserCreateRequest{
				Name:     "J",
				Email:    "john@example.com",
				Password: "password123",
			},
			wantErr: "name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := setupHandler()
			bodyBytes, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			h.Create(rec, req)

			if rec.Code != http.StatusUnprocessableEntity {
				t.Errorf("expected 422 UnprocessableEntity, got %d", rec.Code)
			}

			var resp map[string]any
			json.NewDecoder(rec.Body).Decode(&resp)
			if errObj, ok := resp["error"].(map[string]any); ok {
				if msg, ok := errObj["message"].(string); ok {
					if !containsAny(msg, tt.wantErr) {
						t.Errorf("expected error containing %q, got %q", tt.wantErr, msg)
					}
				}
			}
		})
	}
}

func TestCreate_Success(t *testing.T) {
	h := setupHandler()
	body := domain.UserCreateRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.Create(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("expected 201 Created, got %d", rec.Code)
	}
}

func TestUpdate_ValidationErrors(t *testing.T) {
	tests := []struct {
		name    string
		body    any
		wantErr string
	}{
		{
			name: "invalid email format",
			body: domain.UserUpdateRequest{
				ID:    "550e8400-e29b-41d4-a716-446655440000",
				Email: "not-an-email",
			},
			wantErr: "email",
		},
		{
			name: "name too short",
			body: domain.UserUpdateRequest{
				ID:   "550e8400-e29b-41d4-a716-446655440000",
				Name: "J",
			},
			wantErr: "name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := setupHandler()
			bodyBytes, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPut, "/api/v1/users/"+tt.body.(domain.UserUpdateRequest).ID, bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			h.Update(rec, req)

			if rec.Code != http.StatusUnprocessableEntity {
				t.Errorf("expected 422 UnprocessableEntity, got %d", rec.Code)
			}

			var resp map[string]any
			json.NewDecoder(rec.Body).Decode(&resp)
			if errObj, ok := resp["error"].(map[string]any); ok {
				if msg, ok := errObj["message"].(string); ok {
					if !containsAny(msg, tt.wantErr) {
						t.Errorf("expected error containing %q, got %q", tt.wantErr, msg)
					}
				}
			}
		})
	}
}

func containsAny(s, substr string) bool {
	// simple substring check (production code would use more robust matching)
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
