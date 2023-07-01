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
	ListProductsByCategory(ctx context.Context, categoryID uuid.UUID, limit, offset int) (*domain.ProductList, error)
}

type CategoriesUseCase interface {
	GetCategory(ctx context.Context, id uuid.UUID) (*domain.Category, error)
	InsertCategory(ctx context.Context, userID uuid.UUID, in *domain.Category) (*domain.Category, error)
	DeleteCategory(ctx context.Context, userID, id uuid.UUID) error
}

type OrdersUseCase interface {
	CreateOrder(ctx context.Context, userID, productID uuid.UUID) (*domain.Order, error)
	InsertProductsIntoOrder(ctx context.Context, userID, orderID, productID uuid.UUID) (*domain.Order, error)
	RemoveProductFromOrder(ctx context.Context, userID, orderID, productID uuid.UUID) (*domain.Order, error)
	FinishOrder(ctx context.Context, userID, orderID uuid.UUID)
}

type CheckoutUseCase interface {
	PayOrder(ctx context.Context, orderID, userID uuid.UUID)
}
