package ports

import (
	"context"
	"github.com/SOAT1StackGoLang/tech-challenge/internal/core/domain"
	"github.com/google/uuid"
)

// UsersUseCase Primary actors
type UsersUseCase interface {
	CreateUser(ctx context.Context, name, document, email string) (*domain.User, error)
	ValidateUser(ctx context.Context, document string) (uuid.UUID, error)
	IsUserAdmin(ctx context.Context, id uuid.UUID) (bool, error)
}

type ProductsUseCase interface {
	GetProduct(ctx context.Context, id uuid.UUID) (*domain.Product, error)
	InsertProduct(ctx context.Context, userID uuid.UUID, product *domain.Product) (*domain.Product, error)
	UpdateProduct(ctx context.Context, userID uuid.UUID, product *domain.Product) (*domain.Product, error)
	DeleteProduct(ctx context.Context, userID uuid.UUID, uuid uuid.UUID) error
	ListProductsByCategory(ctx context.Context, category string) ([]*domain.Product, error)
}
