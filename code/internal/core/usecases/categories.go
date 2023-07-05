package usecases

import (
	"context"
	"github.com/SOAT1StackGoLang/tech-challenge/internal/core/domain"
	"github.com/SOAT1StackGoLang/tech-challenge/internal/core/ports"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type categoriesUseCase struct {
	log     *zap.SugaredLogger
	catRepo ports.CategoriesRepository
	userUC  ports.UsersUseCase
}

func NewCategoriesUseCase(logger *zap.SugaredLogger, repo ports.CategoriesRepository, userUC ports.UsersUseCase) ports.CategoriesUseCase {
	return &categoriesUseCase{log: logger, catRepo: repo, userUC: userUC}
}

func (c *categoriesUseCase) GetCategory(ctx context.Context, id uuid.UUID) (*domain.Category, error) {
	out, err := c.catRepo.GetCategoryByID(ctx, id)
	return out, err
}

func (c *categoriesUseCase) InsertCategory(ctx context.Context, userID uuid.UUID, in *domain.Category) (*domain.Category, error) {
	if err := validateIsAdmin(c.log, c.userUC, ctx, userID); err != nil {
		return nil, err
	}
	newCat := domain.NewCategory(uuid.New(), in.CreatedAt, in.Name)

	out, err := c.catRepo.InsertCategory(ctx, newCat)

	return out, err
}

func (c *categoriesUseCase) DeleteCategory(ctx context.Context, userID, id uuid.UUID) error {
	if err := validateIsAdmin(c.log, c.userUC, ctx, userID); err != nil {
		return err
	}

	err := c.catRepo.DeleteCategory(ctx, id)
	return err
}
