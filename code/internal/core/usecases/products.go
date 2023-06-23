package usecases

import (
	"context"
	"errors"
	"github.com/SOAT1StackGoLang/tech-challenge/internal/core/domain"
	"github.com/SOAT1StackGoLang/tech-challenge/internal/core/ports"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var ErrUnauthorized = errors.New("user is not authorize to access this resource")

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
	err := p.validateIsAdmin(ctx, userID)
	if err != nil {
		return nil, err
	}

	out, err := p.productRepo.InsertProduct(ctx, in)
	return out, err
}

func (p productsUseCase) UpdateProduct(ctx context.Context, userID uuid.UUID, in *domain.Product) (*domain.Product, error) {
	err := p.validateIsAdmin(ctx, userID)
	if err != nil {
		return nil, err
	}

	out, err := p.productRepo.UpdateProduct(ctx, in)
	return out, nil
}

func (p productsUseCase) DeleteProduct(ctx context.Context, userID uuid.UUID, prodID uuid.UUID) error {
	err := p.validateIsAdmin(ctx, userID)
	if err != nil {
		return err
	}

	err = p.productRepo.DeleteProduct(ctx, prodID)
	return err
}

func (p productsUseCase) ListProductsByCategory(ctx context.Context, category string) ([]*domain.Product, error) {
	out, err := p.productRepo.ListProductsByCategory(ctx, category)
	return out, err

}

func (p productsUseCase) validateIsAdmin(ctx context.Context, userID uuid.UUID) error {
	admin, err := p.userUC.IsUserAdmin(ctx, userID)
	switch {
	case err != nil:
		p.logger.Errorw(
			"failed checking user is admin",
			zap.String("userID", userID.String()),
			zap.Error(err),
		)
		return err
	case !admin:
		p.logger.Errorw(
			"unauthorized user",
			zap.String("id", userID.String()),
			zap.Error(ErrUnauthorized),
		)
		err = ErrUnauthorized
		return err
	}
	return nil
}
