package domain

import (
	"github.com/google/uuid"
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

func NewUser(code string, name string, email string) *User {
	id := uuid.New()
	return &User{ID: id, Document: code, Name: name, Email: email}
}

type Product struct {
	ID          uuid.UUID
	Category    string
	Name        string
	Description string
	Price       string
}

func NewProduct(category, name, description, price string) *Product {
	return &Product{Category: category, Name: name, Description: description, Price: price}
}

type ProductList struct {
	Products      []*Product
	Limit, Offset int
	Total         int64
}

type Order struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	PaymentID   uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
	Status      string
	ProductsIDs []uuid.UUID
}

func NewOrder(ID uuid.UUID, userID uuid.UUID, paymentID uuid.UUID, createdAt time.Time, updatedAt time.Time, deletedAt time.Time, status string, products []uuid.UUID) *Order {
	return &Order{ID: ID, UserID: userID, PaymentID: paymentID, CreatedAt: createdAt, UpdatedAt: updatedAt, DeletedAt: deletedAt, Status: status, ProductsIDs: products}
}

type Category struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
}

func NewCategory(createdAt time.Time, updatedAt time.Time, name string) *Category {
	return &Category{CreatedAt: createdAt, UpdatedAt: updatedAt, Name: name}
}
