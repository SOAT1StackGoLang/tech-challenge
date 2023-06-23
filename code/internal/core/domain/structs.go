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

type Order struct {
	ID        uuid.UUID
	User      User
	Payment   PaymentInfo
	CreatedAt time.Time
	UpdatedAt time.Time
	Status    OrderStatus
	Items     []Product
}

func NewOrder(user User, payment PaymentInfo, createdAt time.Time, updatedAt time.Time, status OrderStatus, items []Product) *Order {
	return &Order{User: user, Payment: payment, CreatedAt: createdAt, UpdatedAt: updatedAt, Status: status, Items: items}
}

type PaymentInfo struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	OrderID   uuid.UUID
	PaymentID uuid.UUID
	Status    string
	Value     float64
}

func NewPaymentInfo(createdAt time.Time, updatedAt time.Time, orderID uuid.UUID, paymentID uuid.UUID, status string, value float64) *PaymentInfo {
	return &PaymentInfo{CreatedAt: createdAt, UpdatedAt: updatedAt, OrderID: orderID, PaymentID: paymentID, Status: status, Value: value}
}
