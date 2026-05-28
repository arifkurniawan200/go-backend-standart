package middleware

import (
	"context"
	"testing"
)

func TestWithClaimsAndGetClaims(t *testing.T) {
	claims := &JWTClaims{
		Sub:   "user-123",
		Email: "test@example.com",
		Role:  "user",
	}

	ctx := WithClaims(context.Background(), claims)
	got, ok := GetClaims(ctx)

	if !ok {
		t.Fatal("expected claims to be found in context")
	}
	if got.Sub != claims.Sub {
		t.Errorf("Sub = %q, want %q", got.Sub, claims.Sub)
	}
	if got.Email != claims.Email {
		t.Errorf("Email = %q, want %q", got.Email, claims.Email)
	}
	if got.Role != claims.Role {
		t.Errorf("Role = %q, want %q", got.Role, claims.Role)
	}
}

func TestGetClaims_EmptyContext(t *testing.T) {
	_, ok := GetClaims(context.Background())
	if ok {
		t.Error("expected claims not to be found in empty context")
	}
}
