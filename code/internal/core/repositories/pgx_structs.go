package repositories

import (
	"database/sql"
	"github.com/SOAT1StackGoLang/tech-challenge/internal/core/domain"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/type/decimal"
	"time"
)

type User struct {
	ID        uuid.UUID    `gorm:"id,primaryKey"`
	CreatedAt time.Time    `gorm:"created_at"`
	UpdatedAt sql.NullTime `gorm:"updated_at"`
	DeletedAt sql.NullTime `gorm:"deleted_at"`
	Document  string       `gorm:"document"`
	Name      string       `gorm:"name"`
	Email     string       `gorm:"email"`
	IsAdmin   bool         `gorm:"is_admin"`
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

type Category struct {
	ID   uuid.UUID `gorm:"id,primaryKey"`
	Name string    `gorm:"name"`
}

func (c *Category) toDomain() *domain.Category {
	return &domain.Category{
		ID:   c.ID,
		Name: c.Name,
	}
}

func (c *Category) fromDomain(dCat *domain.Category) {
	if c == nil {
		c = &Category{}
	}

	c.ID = dCat.ID
	c.Name = dCat.Name
}

type Product struct {
	ID          uuid.UUID       `gorm:"id,primaryKey"`
	CategoryID  uuid.UUID       `gorm:"category_id"`
	Name        string          `gorm:"name"`
	Image       string          `gorm:"image"`
	Description string          `gorm:"description"`
	Price       decimal.Decimal `gorm:"price"`
}

func (p *Product) toDomain() *domain.Product {
	return &domain.Product{
		ID: p.ID,
		//CategoryID:  p.,
		Name:        p.Name,
		Image:       p.Image,
		Description: p.Description,
		Price:       p.Price.GetValue(),
	}
}
