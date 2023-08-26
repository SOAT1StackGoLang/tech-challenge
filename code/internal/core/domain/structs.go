package domain

import (
	"github.com/SOAT1StackGoLang/tech-challenge/helpers"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

type User struct {
	ID       uuid.UUID
	Document string
	Name     string
	Email    string
}

func NewUser(ID uuid.UUID, document string, name string, email string) *User {
	return &User{ID: ID, Document: document, Name: name, Email: email}
}

type Product struct {
	ID          uuid.UUID
	CategoryID  uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
	Name        string
	Description string
	Price       decimal.Decimal
}

func ParseProductToDomain(
	ID uuid.UUID,
	categoryID uuid.UUID,
	name string,
	description string,
	price string,
) *Product {
	p, err := helpers.ParseDecimalFromString(price)
	if err != nil {
		return nil
	}
	return &Product{
		ID:          ID,
		CategoryID:  categoryID,
		Name:        name,
		Description: description,
		Price:       p,
	}
}

func NewProduct(ID uuid.UUID, categoryID uuid.UUID, name string, description string, price string) *Product {
	p, _ := helpers.ParseDecimalFromString(price)
	return &Product{ID: ID, CategoryID: categoryID, CreatedAt: time.Now(), Name: name, Description: description, Price: p}
}

type ProductList struct {
	Products      []*Product
	Limit, Offset int
	Total         int64
}

type OrderList struct {
	Orders        []*Order
	Limit, Offset int
	Total         int64
}

type ProductsSum struct {
	Products    []uuid.UUID
	RequestedAt time.Time
	Sum         decimal.Decimal
}

type Order struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	PaymentID uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
	Price     decimal.Decimal
	Status    string
	Products  []Product
}

func NewOrder(ID uuid.UUID, userID uuid.UUID, createdAt time.Time, products []Product) *Order {
	return &Order{ID: ID, UserID: userID, CreatedAt: createdAt, Products: products, Status: ORDER_STATUS_OPEN}
}

type Category struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
}

func NewCategory(ID uuid.UUID, createdAt time.Time, name string) *Category {
	return &Category{ID: ID, CreatedAt: createdAt, Name: name}
}

type Payment struct {
	ID        uuid.UUID
	CreatedAt time.Time
	PaidAt    time.Time
	Price     decimal.Decimal
	OrderID   uuid.UUID
}

func NewPayment(ID uuid.UUID, createdAt time.Time, orderID uuid.UUID, price decimal.Decimal) *Payment {
	return &Payment{ID: ID, CreatedAt: createdAt, OrderID: orderID, Price: price}
}

const (
	ORDER_STATUS_UNSET           string = ""
	ORDER_STATUS_OPEN                   = "Aberto"
	ORDER_STATUS_WAITING_PAYMENT        = "Aguardando Pagamento"
	ORDER_STATUS_RECEIVED               = "Recebido"
	ORDER_STATUS_PREPARING              = "Em Preparação"
	ORDER_STATUS_DONE                   = "Pronto"
	ORDER_STATUS_FINISHED               = "Finalizado"
)
