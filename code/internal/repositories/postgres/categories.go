package postgres

import (
	"context"
	"github.com/SOAT1StackGoLang/tech-challenge/internal/core/domain"
	"github.com/SOAT1StackGoLang/tech-challenge/internal/core/ports"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const categoriesTable = "lanchonete_categories"

type categoriesRepositoryImpl struct {
	log *zap.SugaredLogger
	db  *gorm.DB
}

func (c categoriesRepositoryImpl) InsertCategory(ctx context.Context, in *domain.Category) (*domain.Category, error) {
	cat := Category{}
	cat.fromDomain(in)

	if err := c.db.WithContext(ctx).Table(categoriesTable).
		Create(&cat).Error; err != nil {
		c.log.Errorw(
			"db failed inserting category",
			zap.Any("in_category", in),
			zap.Error(err),
		)
		return nil, err
	}

	return cat.toDomain(), nil
}

func (c categoriesRepositoryImpl) GetCategoryByID(ctx context.Context, id uuid.UUID) (*domain.Category, error) {
	cat := Category{}

	if err := c.db.WithContext(ctx).Table(categoriesTable).
		Select("*").Where("id = ?", id).First(&cat).Error; err != nil {
		c.log.Errorw(
			"db failed getting category",
			zap.String("category_id", id.String()),
			zap.Error(err),
		)
		return nil, err
	}

	return cat.toDomain(), nil
}

func (c categoriesRepositoryImpl) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	cat := Category{ID: id}
	if err := c.db.WithContext(ctx).Table(categoriesTable).Delete(&cat).Error; err != nil {
		c.log.Errorw(
			"db failed deleting category",
			zap.Any("category_id", id.String()),
			zap.Error(err),
		)
		return err
	}

	return nil
}

func NewPgxCategoriesRepository(db *gorm.DB, logger *zap.SugaredLogger) ports.CategoriesRepository {
	return &categoriesRepositoryImpl{
		log: logger,
		db:  db,
	}
}
