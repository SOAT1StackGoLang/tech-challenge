package postgres

import (
	"database/sql"
	"github.com/SOAT1StackGoLang/tech-challenge/internal/core/domain"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
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
	ID          uuid.UUID       `gorm:"id,primaryKey" json:"id"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   sql.NullTime    `json:"updated_at"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	CategoryID  uuid.UUID       `json:"category_id"`
	Price       decimal.Decimal `json:"price"`
}

func (p *Product) toDomain() *domain.Product {
	return &domain.Product{
		ID:          p.ID,
		CategoryID:  p.CategoryID,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt.Time,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price.String(),
	}
}

func (p *Product) fromDomain(dProd *domain.Product) {
	if p == nil {
		p = &Product{}
	}

	var decimalValue decimal.Decimal
	var err error
	if dProd.Price != "" {
		decimalValue, err = decimal.NewFromString(dProd.Price)
		if err != nil {
			return
		}
	} else {
		decimalValue = decimal.NewFromInt(0)
	}

	p.ID = dProd.ID
	p.Name = dProd.Name
	p.CategoryID = dProd.CategoryID
	p.Description = dProd.Description
	p.Price = decimalValue
}

type ProductList struct {
	products      []*domain.Product
	limit, offset int
	total         int64
}

type Category struct {
	ID        uuid.UUID `gorm:"id,primaryKey"`
	CreatedAt time.Time
	UpdatedAt sql.NullTime
	Name      string
}

func (c *Category) toDomain() *domain.Category {
	return &domain.Category{
		ID:        c.ID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt.Time,
		Name:      c.Name,
	}
}

func (c *Category) fromDomain(in *domain.Category) {
	if c == nil {
		c = &Category{}
	}

	c.ID = in.ID
	c.CreatedAt = in.CreatedAt
	if !in.UpdatedAt.IsZero() {
		c.UpdatedAt.Time = in.UpdatedAt
		c.UpdatedAt.Valid = true
	}
	c.Name = in.Name

	return
}

type Order struct {
	ID        uuid.UUID `gorm:"id,primaryKey"`
	UserID    uuid.UUID
	PaymentID uuid.UUID
	CreatedAt time.Time
	UpdatedAt sql.NullTime
	DeletedAt sql.NullTime
	Price     decimal.Decimal
	Status    string
	Products  []uuid.UUID `gorm:"type:jsonb"`
}

func (o *Order) newFromDomain(order *domain.Order) {
	if o == nil {
		o = &Order{}
	}
	o = &Order{ID: uuid.New(), UserID: order.UserID, Products: order.ProductsIDs, CreatedAt: time.Now(), Status: "OPEN"}
}

func (o *Order) fromDomain(order *domain.Order) {
	if o == nil {
		o = &Order{}
	}
	o.ID = order.ID
	o.UserID = order.UserID
	o.PaymentID = order.PaymentID
	o.CreatedAt = order.CreatedAt
	o.Status = order.Status
	o.Price = order.Price
	o.Products = order.ProductsIDs

	if !order.UpdatedAt.IsZero() {
		o.UpdatedAt = sql.NullTime{
			Time:  order.UpdatedAt,
			Valid: true,
		}
	}
	if !order.DeletedAt.IsZero() {
		o.DeletedAt = sql.NullTime{
			Time:  order.DeletedAt,
			Valid: true,
		}
	}
}

func (o *Order) toDomain() *domain.Order {
	return &domain.Order{
		ID:          o.ID,
		UserID:      o.UserID,
		PaymentID:   o.PaymentID,
		CreatedAt:   o.CreatedAt,
		UpdatedAt:   o.UpdatedAt.Time,
		DeletedAt:   o.DeletedAt.Time,
		Status:      o.Status,
		Price:       o.Price,
		ProductsIDs: o.Products,
	}
}

type Payment struct {
	ID        uuid.UUID `gorm:"id,primaryKey"`
	CreatedAt time.Time
	OrderID   uuid.UUID
	UserID    uuid.UUID
}

func (p *Payment) fromDomain(dP *domain.Payment) {
	if p == nil {
		p = &Payment{ID: uuid.New(), CreatedAt: time.Now()}
	}
	p.OrderID = dP.OrderID
	p.UserID = dP.UserID
	p.CreatedAt = dP.CreatedAt
}

func (p *Payment) newPayment(userID, orderID uuid.UUID) {
	if p == nil {
		p = &Payment{ID: uuid.New(), CreatedAt: time.Now()}
	}
	p.OrderID = orderID
	p.UserID = userID
	p.CreatedAt = time.Now()
}

func (p *Payment) toDomain() *domain.Payment {
	return &domain.Payment{
		ID:        p.ID,
		CreatedAt: p.CreatedAt,
		OrderID:   p.OrderID,
		UserID:    p.UserID,
	}
}
