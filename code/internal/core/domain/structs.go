package domain

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

type OrderStatus int

const (
	OrderStatusUnset OrderStatus = iota
	OrderStatusOpen
	OrderStatusPaid
	OrderStatusReverted
	OrderStatusAccepted
	OrderStatusPreparing
	OrderStatusInTransit
	OrderStatusDelivered
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
	Price       string
}

func ParseProductToDomain(ID uuid.UUID, categoryID uuid.UUID, createdAt time.Time, updatedAt time.Time, deletedAt time.Time, name string, description string, price string) *Product {
	return &Product{ID: ID, CategoryID: categoryID, CreatedAt: createdAt, UpdatedAt: updatedAt, DeletedAt: deletedAt, Name: name, Description: description, Price: price}
}

func NewProduct(ID uuid.UUID, categoryID uuid.UUID, name string, description string, price string) *Product {
	return &Product{ID: ID, CategoryID: categoryID, CreatedAt: time.Now(), Name: name, Description: description, Price: price}
}

type ProductList struct {
	Products      []*Product
	Limit, Offset int
	Total         int64
}

type ProductsSum struct {
	Products    []uuid.UUID
	RequestedAt time.Time
	Sum         decimal.Decimal
}

type Order struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	PaymentID   uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
	Price       string
	Status      string
	ProductsIDs []uuid.UUID
}

func NewOrder(ID uuid.UUID, userID uuid.UUID, paymentID uuid.UUID, createdAt time.Time, products []uuid.UUID) *Order {
	return &Order{ID: ID, UserID: userID, PaymentID: paymentID, CreatedAt: createdAt, ProductsIDs: products}
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
	OrderID   uuid.UUID
	UserID    uuid.UUID
}

func NewPayment(ID uuid.UUID, createdAt time.Time, orderID uuid.UUID, userID uuid.UUID) *Payment {
	return &Payment{ID: ID, CreatedAt: createdAt, OrderID: orderID, UserID: userID}
}
