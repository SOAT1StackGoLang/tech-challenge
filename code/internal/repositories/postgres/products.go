package postgres

import (
	"context"
	"database/sql"
	"github.com/SOAT1StackGoLang/tech-challenge/internal/core/domain"
	"github.com/SOAT1StackGoLang/tech-challenge/internal/core/ports"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

const productsTable = "lanchonete_products"

type productsRepositoryImpl struct {
	log *zap.SugaredLogger
	db  *gorm.DB
}

func NewPgxProductsRepository(db *gorm.DB, logger *zap.SugaredLogger) ports.ProductsRepository {
	return &productsRepositoryImpl{
		log: logger,
		db:  db,
	}
}

func (p productsRepositoryImpl) GetProduct(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	out := Product{}

	if err := p.db.WithContext(ctx).Table(productsTable).
		Select("*").Where("id = ?", id).First(&out).Error; err != nil {
		p.log.Errorw(
			"db failed getting product",
			zap.String("id", id.String()),
			zap.Error(err),
		)
		return nil, err
	}

	return out.toDomain(), nil
}

func (p productsRepositoryImpl) InsertProduct(ctx context.Context, in *domain.Product) (*domain.Product, error) {
	product := Product{}
	product.fromDomain(in)

	if err := p.db.WithContext(ctx).Table(productsTable).Create(&product).Error; err != nil {
		p.log.Errorw(
			"db failed inserting product",
			zap.Any("in_product", in),
			zap.Error(err),
		)
		return nil, err
	}

	return product.toDomain(), nil

}

func (p productsRepositoryImpl) UpdateProduct(ctx context.Context, in *domain.Product) (*domain.Product, error) {
	product := Product{}
	product.fromDomain(in)
	product.UpdatedAt = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	if err := p.db.WithContext(ctx).Table(productsTable).Updates(&product).Where("id = ?", in.ID).Error; err != nil {
		p.log.Errorw(
			"db failed updating product",
			zap.Any("in_product", in),
			zap.Error(err),
		)
		return nil, err
	}

	return product.toDomain(), nil
}

func (p productsRepositoryImpl) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	product := Product{ID: id}
	if err := p.db.WithContext(ctx).Table(productsTable).Delete(&product).Error; err != nil {
		p.log.Errorw(
			"failed deleting product",
			zap.String("product_id", id.String()),
			zap.Error(err),
		)
		return err
	}

	return nil
}

func (p productsRepositoryImpl) ListProductsByCategory(ctx context.Context, categoryID uuid.UUID, limit, offset int) (*domain.ProductList, error) {
	var products []Product
	var total int64

	err := p.db.WithContext(ctx).Table(productsTable).
		Where("category_id = ?", categoryID).Limit(limit).Offset(offset).Order("name ASC").Find(&products).Error
	if err != nil {
		p.log.Errorw(
			"failed listing products",
			zap.String("category", categoryID.String()),
			zap.Error(err),
		)
		return nil, err
	}

	_ = p.db.WithContext(ctx).Table(productsTable).
		Where("category_id = ?").Count(&total)

	var pList *domain.ProductList
	var out []*domain.Product
	for _, v := range products {
		out = append(out, v.toDomain())
	}

	pList.Products = out
	pList.Total = total
	pList.Limit = limit
	pList.Offset = offset

	return pList, err
}
