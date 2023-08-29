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

func NewPgxCategoriesRepository(db *gorm.DB, logger *zap.SugaredLogger) ports.CategoriesRepository {
	return &categoriesRepositoryImpl{
		log: logger,
		db:  db,
	}
}

func (c categoriesRepositoryImpl) ListCategories(ctx context.Context, limit int, offset int) (*domain.CategoryList, error) {
	var total int64
	var savedCats []Category

	var err error
	if err = c.db.WithContext(ctx).Table(categoriesTable).
		Limit(limit).
		Offset(offset).
		Count(&total).
		Scan(&savedCats).Error; err != nil {
		c.log.Errorw(
			"failed listing categories",
			zap.Error(err),
		)
		return nil, err
	}

	out := &domain.CategoryList{}
	outList := make([]*domain.Category, 0, total)

	for _, c := range savedCats {
		outList = append(outList, c.toDomain())
	}

	out.Categories = outList
	out.Total = total
	out.Limit = limit
	out.Offset = offset

	return out, err
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
