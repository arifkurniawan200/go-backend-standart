package usecase

import (
	"context"
	"errors"

	"github.com/arifkurniawan200/go-backend-standart/internal/domain"
	"github.com/arifkurniawan200/go-backend-standart/internal/repository"
)

// Common errors
var (
	ErrInvalidInput = errors.New("invalid input")
)

// UserUsecase handles user business logic
type UserUsecase struct {
	repo repository.UserRepository
}

// NewUserUsecase creates a new user usecase
func NewUserUsecase(repo repository.UserRepository) *UserUsecase {
	return &UserUsecase{repo: repo}
}

// Create creates a new user
func (uc *UserUsecase) Create(ctx context.Context, req *domain.UserCreateRequest) (*domain.User, error) {
	if req.Name == "" || req.Email == "" || req.Password == "" {
		return nil, ErrInvalidInput
	}

	user := &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password, // In production: hash this!
	}

	if err := uc.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetByID retrieves a user by ID
func (uc *UserUsecase) GetByID(ctx context.Context, id string) (*domain.User, error) {
	return uc.repo.FindByID(ctx, id)
}

// GetByEmail retrieves a user by email
func (uc *UserUsecase) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return uc.repo.FindByEmail(ctx, email)
}

// Update updates an existing user
func (uc *UserUsecase) Update(ctx context.Context, req *domain.UserUpdateRequest) (*domain.User, error) {
	if req.ID == "" {
		return nil, ErrInvalidInput
	}

	existingUser, err := uc.repo.FindByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		existingUser.Name = req.Name
	}
	if req.Email != "" {
		existingUser.Email = req.Email
	}

	if err := uc.repo.Update(ctx, existingUser); err != nil {
		return nil, err
	}

	return existingUser, nil
}

// Delete removes a user
func (uc *UserUsecase) Delete(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}

// List returns all users with pagination
func (uc *UserUsecase) List(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	return uc.repo.List(ctx, limit, offset)
}
