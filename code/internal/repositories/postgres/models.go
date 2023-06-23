package postgres

import (
	"database/sql"
	"github.com/SOAT1StackGoLang/tech-challenge/internal/core/domain"
	"github.com/shopspring/decimal"

	"github.com/google/uuid"
	"time"
)

type User struct {
	ID        uuid.UUID `gorm:"id,primaryKey"`
	CreatedAt time.Time
	UpdatedAt sql.NullTime
	DeletedAt sql.NullTime
	Document  string
	Name      string
	Email     string
	IsAdmin   bool
}

func (u *User) toDomain() *domain.User {
	return &domain.User{ID: u.ID, Document: u.Document, Name: u.Name, Email: u.Email}
}

func (u *User) fromDomain(dUser *domain.User) {
	if u == nil {
		u = &User{}
	}

	u.ID = dUser.ID
	u.Document = dUser.Document
	u.Name = dUser.Name
	u.Email = dUser.Email
}

type Product struct {
	ID          uuid.UUID `gorm:"id,primaryKey"`
	CreatedAt   time.Time
	Name        string
	Description string
	Category    string
	Price       decimal.Decimal
}

func (p *Product) toDomain() *domain.Product {
	return &domain.Product{
		ID:          p.ID,
		Category:    p.Category,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price.String(),
	}
}

func (p *Product) fromDomain(dProd *domain.Product) {
	if p == nil {
		p = &Product{}
	}

	decimalValue, err := decimal.NewFromString(dProd.Price)
	if err != nil {
		return
	}

	p.ID = dProd.ID
	p.Name = dProd.Name
	p.Category = dProd.Category
	p.Description = dProd.Description
	p.Price = decimalValue
}
