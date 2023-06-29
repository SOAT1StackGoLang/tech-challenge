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

type ProductsRepository interface {
	GetProduct(ctx context.Context, id uuid.UUID) (*domain.Product, error)
	InsertProduct(ctx context.Context, product *domain.Product) (*domain.Product, error)
	UpdateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error)
	DeleteProduct(ctx context.Context, uuid uuid.UUID) error
	ListProductsByCategory(ctx context.Context, categoryID uuid.UUID, limit, offset int) (*domain.ProductList, error)
}

type CategoriesRepository interface {
	InsertCategory(ctx context.Context, in *domain.Category) (*domain.Category, error)
	GetCategoryByID(ctx context.Context, id uuid.UUID) (*domain.Category, error)
	DeleteCategory(ctx context.Context, id uuid.UUID) error
}

type OrdersRepository interface {
	CreateOrder(ctx context.Context, userID, productID uuid.UUID) (*domain.Order, error)
	InsertProductIntoOrder(ctx context.Context, userID, orderID, productID uuid.UUID) (*domain.Order, error)
	RemoveProductFromOrder(ctx context.Context, userID, orderID, productID uuid.UUID) (*domain.Order, error)
	DeleteOrder(ctx context.Context, userID, orderID uuid.UUID) error
}