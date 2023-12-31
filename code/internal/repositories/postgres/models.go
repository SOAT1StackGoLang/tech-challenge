package postgres

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/SOAT1StackGoLang/tech-challenge/internal/core/domain"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
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

type OrderProduct struct {
	ID          uuid.UUID       `gorm:"id,primaryKey" json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	CategoryID  uuid.UUID       `json:"category_id"`
	Price       decimal.Decimal `json:"price"`
}

func (oP *OrderProduct) toDomain() *domain.Product {
	return &domain.Product{
		ID:          oP.ID,
		Name:        oP.Name,
		Description: oP.Description,
		CategoryID:  oP.CategoryID,
		Price:       oP.Price,
	}
}

func (oP *OrderProduct) fromDomain(dP *domain.Product) {
	oP.ID = dP.ID
	oP.Name = dP.Name
	oP.Description = dP.Description
	oP.CategoryID = dP.CategoryID
	oP.Price = dP.Price
}

func (p *Product) toDomain() *domain.Product {
	return &domain.Product{
		ID:          p.ID,
		CategoryID:  p.CategoryID,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt.Time,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
	}
}

func (p *Product) fromDomain(dProd *domain.Product) {
	if p == nil {
		p = &Product{}
	}

	p.ID = dProd.ID
	p.Name = dProd.Name
	p.CategoryID = dProd.CategoryID
	p.Description = dProd.Description
	p.Price = dProd.Price
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
	Status    OrderStatus
	Products  json.RawMessage `json:"products" gorm:"type:jsonb"`
}

func (o *Order) fromDomain(order *domain.Order) {
	if o == nil {
		panic("order is nil")
	}
	o.ID = order.ID
	o.UserID = order.UserID
	o.PaymentID = order.PaymentID
	o.CreatedAt = order.CreatedAt
	o.Price = order.Price

	var products []OrderProduct
	for _, p := range order.Products {
		oP := OrderProduct{}
		oP.fromDomain(&p)
		products = append(products, oP)
	}

	productsJSON, err := json.Marshal(products)
	if err != nil {
		//TODO handle properly
		panic("failed to marshal products")
	}

	o.Products = productsJSON

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

	var oS OrderStatus
	o.Status = oS.fromDomain(order.Status)
}

func (o *Order) toDomain() *domain.Order {
	var products []Product
	err := json.Unmarshal(o.Products, &products)
	if err != nil {
		// TODO handle properly
		panic("failed to unmarshal products")
	}

	var outProducts []domain.Product
	for _, v := range products {
		outProducts = append(outProducts, *v.toDomain())
	}

	return &domain.Order{
		ID:        o.ID,
		UserID:    o.UserID,
		PaymentID: o.PaymentID,
		CreatedAt: o.CreatedAt,
		UpdatedAt: o.UpdatedAt.Time,
		DeletedAt: o.DeletedAt.Time,
		Status:    o.Status.toDomain(),
		Price:     o.Price,
		Products:  outProducts,
	}
}

func (o *Order) extractProductsIDsFromJSON() []uuid.UUID {
	var products []string
	err := json.Unmarshal(o.Products, &products)
	if err != nil {
		// Handle error
	}

	// Convert the slice of strings to a slice of uuid.UUID values
	productIDs := make([]uuid.UUID, len(products))
	for i, id := range products {
		productIDs[i], err = uuid.Parse(id)
		if err != nil {
			// Handle error
		}
	}

	return productIDs
}

type OrderStatus int

const (
	ORDER_STATUS_UNSET OrderStatus = iota
	ORDER_STATUS_OPEN
	ORDER_STATUS_WAITING_PAYMENT
	ORDER_STATUS_RECEIVED
	ORDER_STATUS_PREPARING
	ORDER_STATUS_DONE
	ORDER_STATUS_FINISHED
	ORDER_STATUS_CANCELED
)

func (oS OrderStatus) toDomain() domain.OrderStatus {
	switch oS {
	case ORDER_STATUS_OPEN:
		return domain.ORDER_STATUS_OPEN
	case ORDER_STATUS_WAITING_PAYMENT:
		return domain.ORDER_STATUS_WAITING_PAYMENT
	case ORDER_STATUS_RECEIVED:
		return domain.ORDER_STATUS_RECEIVED
	case ORDER_STATUS_PREPARING:
		return domain.ORDER_STATUS_PREPARING
	case ORDER_STATUS_DONE:
		return domain.ORDER_STATUS_DONE
	case ORDER_STATUS_FINISHED:
		return domain.ORDER_STATUS_FINISHED
	case ORDER_STATUS_CANCELED:
		return domain.ORDER_STATUS_CANCELED
	}
	return domain.ORDER_STATUS_UNSET
}

func (oS OrderStatus) fromDomain(status domain.OrderStatus) OrderStatus {
	switch status {
	case domain.ORDER_STATUS_OPEN:
		return ORDER_STATUS_OPEN
	case domain.ORDER_STATUS_WAITING_PAYMENT:
		return ORDER_STATUS_WAITING_PAYMENT
	case domain.ORDER_STATUS_RECEIVED:
		return ORDER_STATUS_RECEIVED
	case domain.ORDER_STATUS_PREPARING:
		return ORDER_STATUS_PREPARING
	case domain.ORDER_STATUS_DONE:
		return ORDER_STATUS_DONE
	case domain.ORDER_STATUS_FINISHED:
		return ORDER_STATUS_FINISHED
	case domain.ORDER_STATUS_CANCELED:
		return ORDER_STATUS_CANCELED
	}
	return ORDER_STATUS_UNSET
}

type Payment struct {
	ID        uuid.UUID `gorm:"id,primaryKey"`
	CreatedAt time.Time
	UpdatedAt sql.NullTime
	Value     decimal.Decimal `json:"value"`
	OrderID   uuid.UUID
	Status    PaymentStatus
}

type PaymentStatus string

const (
	PAYMENT_SATUS_OPEN      PaymentStatus = "Aberto"
	PAYMENT_STATUS_APPROVED               = "Aprovado"
	PAYMENT_STATUS_REFUSED                = "Recusado"
)

func (pS PaymentStatus) toDomain() domain.PaymentStatus {
	switch pS {
	case PAYMENT_SATUS_OPEN:
		return domain.PAYMENT_STATUS_OPEN
	case PAYMENT_STATUS_APPROVED:
		return domain.PAYMENT_STATUS_APPROVED
	}
	return domain.PAYMENT_SATUS_REFUSED
}

func (pS PaymentStatus) fromDomain(in domain.PaymentStatus) PaymentStatus {
	switch in {
	case domain.PAYMENT_STATUS_APPROVED:
		return PAYMENT_STATUS_APPROVED
	case domain.PAYMENT_SATUS_REFUSED:
		return PAYMENT_STATUS_REFUSED

	}
	return PAYMENT_SATUS_OPEN
}

func (p *Payment) fromDomain(dP *domain.Payment) {
	p.ID = dP.ID
	p.OrderID = dP.OrderID
	p.CreatedAt = dP.CreatedAt
	p.Value = dP.Price

	if !dP.UpdatedAt.IsZero() {
		p.UpdatedAt = sql.NullTime{
			Time:  dP.UpdatedAt,
			Valid: true,
		}
	} else {
		p.UpdatedAt = sql.NullTime{Valid: false}
	}

	var pS PaymentStatus
	pS = pS.fromDomain(dP.Status)
	p.Status = pS
}

func (p *Payment) toDomain() *domain.Payment {
	var dS domain.PaymentStatus
	dS = p.Status.toDomain()
	return &domain.Payment{
		ID:        p.ID,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt.Time,
		Price:     p.Value,
		OrderID:   p.OrderID,
		Status:    dS,
	}
}
