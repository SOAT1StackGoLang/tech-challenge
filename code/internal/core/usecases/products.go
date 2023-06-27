package usecases

import (
	"context"
	"github.com/SOAT1StackGoLang/tech-challenge/internal/core/domain"
	"github.com/SOAT1StackGoLang/tech-challenge/internal/core/ports"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type productsUseCase struct {
	logger      *zap.SugaredLogger
	productRepo ports.ProductsRepository
	userUC      ports.UsersUseCase
}

func NewProductsUseCase(repository ports.ProductsRepository, userUseCase ports.UsersUseCase, logger *zap.SugaredLogger) ports.ProductsUseCase {
	return &productsUseCase{
		logger:      logger,
		productRepo: repository,
		userUC:      userUseCase,
	}
}

func (p productsUseCase) GetProduct(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	product, err := p.productRepo.GetProduct(ctx, id)
	return product, err
}

func (p productsUseCase) InsertProduct(ctx context.Context, userID uuid.UUID, in *domain.Product) (*domain.Product, error) {
	err := validateIsAdmin(p.logger, p.userUC, ctx, userID)
	if err != nil {
		return nil, err
	}

	out, err := p.productRepo.InsertProduct(ctx, in)
	return out, err
}

func (p productsUseCase) UpdateProduct(ctx context.Context, userID uuid.UUID, in *domain.Product) (*domain.Product, error) {
	err := validateIsAdmin(p.logger, p.userUC, ctx, userID)
	if err != nil {
		return nil, err
	}

	out, err := p.productRepo.UpdateProduct(ctx, in)
	return out, nil
}

func (p productsUseCase) DeleteProduct(ctx context.Context, userID uuid.UUID, prodID uuid.UUID) error {
	err := validateIsAdmin(p.logger, p.userUC, ctx, userID)
	if err != nil {
		return err
	}

	err = p.productRepo.DeleteProduct(ctx, prodID)
	return err
}

func (p productsUseCase) ListProductsByCategory(ctx context.Context, categoryID uuid.UUID, limit, offset int) (*domain.ProductList, error) {
	out, err := p.productRepo.ListProductsByCategory(ctx, categoryID, limit, offset)
	return out, err

}
