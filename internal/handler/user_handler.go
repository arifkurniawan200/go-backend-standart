package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/arifkurniawan200/go-backend-standart/internal/domain"
	"github.com/arifkurniawan200/go-backend-standart/internal/repository"
	"github.com/arifkurniawan200/go-backend-standart/internal/usecase"
	"github.com/arifkurniawan200/go-backend-standart/pkg/response"
	"go.uber.org/zap"
)

// UserHandler handles HTTP requests for users
type UserHandler struct {
	uc    *usecase.UserUsecase
	log   *zap.Logger
}

// NewUserHandler creates a new user handler
func NewUserHandler(uc *usecase.UserUsecase, log *zap.Logger) *UserHandler {
	return &UserHandler{
		uc:  uc,
		log: log,
	}
}

// RegisterRoutes registers all user routes
func (h *UserHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/users", h.Create)
	mux.HandleFunc("GET /api/v1/users/{id}", h.GetByID)
	mux.HandleFunc("GET /api/v1/users", h.List)
	mux.HandleFunc("PUT /api/v1/users/{id}", h.Update)
	mux.HandleFunc("DELETE /api/v1/users/{id}", h.Delete)
}

// Create handles POST /api/v1/users
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var req domain.UserCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Warn("Invalid request body", "error", err)
		response.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := h.uc.Create(ctx, &req)
	if err != nil {
		h.log.Error("Failed to create user", "error", err)
		response.WriteError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	h.log.Info("User created", "user_id", user.ID)
	response.WriteJSON(w, http.StatusCreated, user)
}

// GetByID handles GET /api/v1/users/{id}
func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	id := r.PathValue("id")
	if id == "" {
		response.WriteError(w, http.StatusBadRequest, "Missing user ID")
		return
	}

	user, err := h.uc.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			response.WriteError(w, http.StatusNotFound, "User not found")
			return
		}
		h.log.Error("Failed to get user", "error", err, "id", id)
		response.WriteError(w, http.StatusInternalServerError, "Failed to get user")
		return
	}

	response.WriteJSON(w, http.StatusOK, user)
}

// List handles GET /api/v1/users
func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	limit := 10
	offset := 0

	users, err := h.uc.List(ctx, limit, offset)
	if err != nil {
		h.log.Error("Failed to list users", "error", err)
		response.WriteError(w, http.StatusInternalServerError, "Failed to list users")
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"data":   users,
		"limit":  limit,
		"offset": offset,
	})
}

// Update handles PUT /api/v1/users/{id}
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	id := r.PathValue("id")
	if id == "" {
		response.WriteError(w, http.StatusBadRequest, "Missing user ID")
		return
	}

	var req domain.UserUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Warn("Invalid request body", "error", err)
		response.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	req.ID = id

	user, err := h.uc.Update(ctx, &req)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			response.WriteError(w, http.StatusNotFound, "User not found")
			return
		}
		h.log.Error("Failed to update user", "error", err, "id", id)
		response.WriteError(w, http.StatusInternalServerError, "Failed to update user")
		return
	}

	h.log.Info("User updated", "user_id", user.ID)
	response.WriteJSON(w, http.StatusOK, user)
}

// Delete handles DELETE /api/v1/users/{id}
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	id := r.PathValue("id")
	if id == "" {
		response.WriteError(w, http.StatusBadRequest, "Missing user ID")
		return
	}

	if err := h.uc.Delete(ctx, id); err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			response.WriteError(w, http.StatusNotFound, "User not found")
			return
		}
		h.log.Error("Failed to delete user", "error", err, "id", id)
		response.WriteError(w, http.StatusInternalServerError, "Failed to delete user")
		return
	}

	h.log.Info("User deleted", "user_id", id)
	response.WriteJSON(w, http.StatusOK, map[string]string{"message": "User deleted successfully"})
}

// HealthCheck for demonstration
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
