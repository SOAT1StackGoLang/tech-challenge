package ports

import (
	"context"
	"github.com/SOAT1StackGoLang/tech-challenge/internal/core/domain"
	"github.com/google/uuid"
)

// UsersRepository Secondary actors
type UsersRepository interface {
	GetUser(ctx context.Context, id uuid.UUID) (*domain.User, error)
	InsertUser(ctx context.Context, user *domain.User) error
	ValidateUser(ctx context.Context, document string) (uuid.UUID, error)
	IsUserAdmin(ctx context.Context, id uuid.UUID) (bool, error)
}
