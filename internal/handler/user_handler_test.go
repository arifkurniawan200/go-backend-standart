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
func (m *mockUserRepo) FindByID(ctx context.Context, id string) (*domain.User, error)    { return nil, nil }
func (m *mockUserRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) { return nil, nil }
func (m *mockUserRepo) Update(ctx context.Context, user *domain.User) error              { return nil }
func (m *mockUserRepo) Delete(ctx context.Context, id string) error                       { return nil }
func (m *mockUserRepo) List(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	return []*domain.User{
		{ID: "1", Name: "Alice", Email: "alice@example.com"},
		{ID: "2", Name: "Bob", Email: "bob@example.com"},
	}, nil
}

func setupHandler() *UserHandler {
	repo := newMockUserRepo()
	uc := usecase.NewUserUsecase(repo)
	log := zaptest.NewLogger(nil)
	return NewUserHandler(uc, log)
}

// ============================================================
// Validator wiring tests (PR #3)
// ============================================================

func TestCreate_InvalidInput(t *testing.T) {
	h := setupHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader([]byte("not json")))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.Create(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestCreate_ValidationErrors(t *testing.T) {
	tests := []struct {
		name    string
		body    any
		wantErr string
	}{
		{name: "missing name", body: domain.UserCreateRequest{Email: "test@ex.com", Password: "pwd123456"}, wantErr: "name"},
		{name: "empty email", body: domain.UserCreateRequest{Name: "John", Email: "", Password: "pwd123456"}, wantErr: "email"},
		{name: "bad email", body: domain.UserCreateRequest{Name: "John", Email: "bad", Password: "pwd123456"}, wantErr: "email"},
		{name: "short pwd", body: domain.UserCreateRequest{Name: "John", Email: "j@ex.com", Password: "sh"}, wantErr: "password"},
		{name: "short name", body: domain.UserCreateRequest{Name: "J", Email: "j@ex.com", Password: "pwd123456"}, wantErr: "name"},
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
				t.Errorf("expected 422, got %d", rec.Code)
			}
			var resp map[string]any
			json.NewDecoder(rec.Body).Decode(&resp)
			if errObj, ok := resp["error"].(map[string]any); ok {
				if msg, ok := errObj["message"].(string); ok {
					if !containsAny(msg, tt.wantErr) {
						t.Errorf("error should contain %q, got %q", tt.wantErr, msg)
					}
				}
			}
		})
	}
}

func TestUpdate_ValidationErrors(t *testing.T) {
	tests := []struct {
		name    string
		body    any
		wantErr string
	}{
		{name: "bad email", body: domain.UserUpdateRequest{ID: "550e8400-e29b-41d4-a716-446655440000", Email: "bad"}, wantErr: "email"},
		{name: "short name", body: domain.UserUpdateRequest{ID: "550e8400-e29b-41d4-a716-446655440000", Name: "J"}, wantErr: "name"},
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
				t.Errorf("expected 422, got %d", rec.Code)
			}
		})
	}
}

func containsAny(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// ============================================================
// Query params tests (PR #4)
// ============================================================

func TestList_DefaultPagination(t *testing.T) {
	h := setupHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
	rec := httptest.NewRecorder()
	h.List(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
	var resp map[string]any
	json.NewDecoder(rec.Body).Decode(&resp)
	if limit, _ := resp["limit"].(float64); limit != 10 {
		t.Errorf("expected default limit 10, got %v", limit)
	}
	if offset, _ := resp["offset"].(float64); offset != 0 {
		t.Errorf("expected default offset 0, got %v", offset)
	}
}

func TestList_CustomPagination(t *testing.T) {
	h := setupHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users?limit=5&offset=20", nil)
	rec := httptest.NewRecorder()
	h.List(rec, req)
	var resp map[string]any
	json.NewDecoder(rec.Body).Decode(&resp)
	if limit, _ := resp["limit"].(float64); limit != 5 {
		t.Errorf("expected limit 5, got %v", limit)
	}
	if offset, _ := resp["offset"].(float64); offset != 20 {
		t.Errorf("expected offset 20, got %v", offset)
	}
}

func TestList_LimitExceedsMax(t *testing.T) {
	h := setupHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users?limit=200", nil)
	rec := httptest.NewRecorder()
	h.List(rec, req)
	var resp map[string]any
	json.NewDecoder(rec.Body).Decode(&resp)
	if limit, _ := resp["limit"].(float64); limit != 100 {
		t.Errorf("expected limit capped at 100, got %v", limit)
	}
}

func TestList_InvalidQueryDefaults(t *testing.T) {
	h := setupHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users?limit=abc&offset=-5", nil)
	rec := httptest.NewRecorder()
	h.List(rec, req)
	var resp map[string]any
	json.NewDecoder(rec.Body).Decode(&resp)
	if limit, _ := resp["limit"].(float64); limit != 10 {
		t.Errorf("expected default limit 10, got %v", limit)
	}
	if offset, _ := resp["offset"].(float64); offset != 0 {
		t.Errorf("expected default offset 0, got %v", offset)
	}
}

// ============================================================
// Basic handler tests (PR #5)
// ============================================================

func TestHandler_Create(t *testing.T) {
	h := setupHandler()
	body := domain.UserCreateRequest{Name: "John", Email: "john@ex.com", Password: "pwd123456"}
	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.Create(rec, req)
	if rec.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", rec.Code)
	}
	var user domain.User
	json.NewDecoder(rec.Body).Decode(&user)
	if user.ID == "" {
		t.Error("expected user ID in response")
	}
}

func TestHandler_GetByID(t *testing.T) {
	h := setupHandler()
	body := domain.UserCreateRequest{Name: "A", Email: "a@ex.com", Password: "pwd123456"}
	b, _ := json.Marshal(body)
	cr := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(b))
	cr.Header.Set("Content-Type", "application/json")
	crr := httptest.NewRecorder()
	h.Create(crr, cr)
	var created domain.User
	json.NewDecoder(crr.Body).Decode(&created)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/"+created.ID, nil)
	rec := httptest.NewRecorder()
	h.GetByID(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestHandler_GetByID_MissingID(t *testing.T) {
	h := setupHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/", nil)
	rec := httptest.NewRecorder()
	h.GetByID(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestHandler_List(t *testing.T) {
	h := setupHandler()
	for i := 0; i < 3; i++ {
		b, _ := json.Marshal(domain.UserCreateRequest{Name: "U", Email: "u@ex.com", Password: "pwd123456"})
		r := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(b))
		r.Header.Set("Content-Type", "application/json")
		rc := httptest.NewRecorder()
		h.Create(rc, r)
	}
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
	rec := httptest.NewRecorder()
	h.List(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestHandler_Update(t *testing.T) {
	h := setupHandler()
	b, _ := json.Marshal(domain.UserCreateRequest{Name: "A", Email: "a@ex.com", Password: "pwd123456"})
	cr := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(b))
	cr.Header.Set("Content-Type", "application/json")
	crr := httptest.NewRecorder()
	h.Create(crr, cr)
	var created domain.User
	json.NewDecoder(crr.Body).Decode(&created)

	ub, _ := json.Marshal(domain.UserUpdateRequest{ID: created.ID, Name: "Updated"})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/users/"+created.ID, bytes.NewReader(ub))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.Update(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestHandler_Update_MissingID(t *testing.T) {
	h := setupHandler()
	ub, _ := json.Marshal(domain.UserUpdateRequest{Name: "A"})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/users/", bytes.NewReader(ub))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.Update(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestHandler_Delete(t *testing.T) {
	h := setupHandler()
	b, _ := json.Marshal(domain.UserCreateRequest{Name: "A", Email: "a@ex.com", Password: "pwd123456"})
	cr := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(b))
	cr.Header.Set("Content-Type", "application/json")
	crr := httptest.NewRecorder()
	h.Create(crr, cr)
	var created domain.User
	json.NewDecoder(crr.Body).Decode(&created)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/users/"+created.ID, nil)
	rec := httptest.NewRecorder()
	h.Delete(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestHandler_Delete_MissingID(t *testing.T) {
	h := setupHandler()
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/users/", nil)
	rec := httptest.NewRecorder()
	h.Delete(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestHandler_HealthCheck(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	HealthCheck(rec, req)
	if rec.Code != http.StatusOK || rec.Body.String() != "OK" {
		t.Errorf("expected 200 OK, got %d %s", rec.Code, rec.Body.String())
	}
}
