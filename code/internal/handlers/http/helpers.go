package http

import (
	"github.com/SOAT1StackGoLang/tech-challenge/helpers"
	"github.com/SOAT1StackGoLang/tech-challenge/internal/core/domain"
	"github.com/google/uuid"
	"time"
)

func (q *QueryStruct) parseToUuid() (uuid.UUID, uuid.UUID, error) {
	id, err := uuid.Parse(q.ID)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}
	uid, err := uuid.Parse(q.UserID)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}
	return id, uid, err
}

func (c *Category) toDomain() *domain.Category {
	if c == nil {
		c = &Category{}
	}

	return &domain.Category{
		ID:        helpers.SafeUUIDFromString(c.ID),
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		Name:      c.Name,
	}
}

func (c *Category) fromDomain(cat *domain.Category) {
	if c == nil {
		c = &Category{}
	}

	c.ID = cat.ID.String()
	c.Name = cat.Name
	c.CreatedAt = cat.CreatedAt
	c.UpdatedAt = cat.UpdatedAt
}

func (o *Order) fromDomain(order *domain.Order) {
	if o == nil {
		o = &Order{}
	}

	o.ID = order.ID

	var p string
	if !order.Price.IsZero() {
		p = helpers.ParseDecimalToString(order.Price)
		o.Price = p
	}

	var products []Product
	for _, dP := range order.Products {
		p := Product{}
		p.fromDomain(&dP)
		products = append(products, p)
	}

	o.Products = products

	if order.PaymentID != uuid.Nil {
		o.PaymentID = order.PaymentID.String()
	}
	o.CreatedAt = order.CreatedAt.Format(time.RFC3339)
	if !order.UpdatedAt.IsZero() {
		o.UpdatedAt = order.UpdatedAt.Format(time.RFC3339)
	} else {
		o.UpdatedAt = ""
	}
	if !order.DeletedAt.IsZero() {
		o.DeletedAt = order.DeletedAt.Format(time.RFC3339)
	} else {
		o.DeletedAt = ""
	}

	orderStatus := new(OrderStatus)
	o.Status = orderStatus.fromDomain(order.Status)
}

func (oS *OrderStatus) fromDomain(status domain.OrderStatus) OrderStatus {
	switch status {
	case domain.ORDER_STATUS_UNSET:
		return ORDER_STATUS_UNSET
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
	}
	return ORDER_STATUS_CANCELED
}

func stringToDomainStatus(status string) domain.OrderStatus {
	switch status {
	case ORDER_STATUS_RECEIVED:
		return domain.ORDER_STATUS_RECEIVED
	case ORDER_STATUS_PREPARING:
		return domain.ORDER_STATUS_PREPARING
	case ORDER_STATUS_DONE:
		return domain.ORDER_STATUS_DONE
	case ORDER_STATUS_FINISHED:
		return domain.ORDER_STATUS_FINISHED
	}
	return domain.ORDER_STATUS_UNSET
}

func (pN *PaymentNotification) toDomain() *domain.PaymentStatusNotification {
	var pS PaymentStatus
	pS = pS.fromRequest(pN.Approved)
	return &domain.PaymentStatusNotification{
		PaymentID: helpers.SafeUUIDFromString(pN.PaymentID),
		OrderID:   helpers.SafeUUIDFromString(pN.OrderID),
		Status:    pS.toDomain(),
	}
}

func (p *Payment) fromDomain(payment *domain.Payment) {
	p.ID = payment.ID.String()
	p.OrderID = payment.OrderID.String()
	p.CreatedAt = payment.CreatedAt.Format(time.RFC3339)

	if !payment.UpdatedAt.IsZero() {
		p.UpdatedAt = payment.UpdatedAt.Format(time.RFC3339)
	}
	var price string
	if !payment.Price.IsZero() {
		price = helpers.ParseDecimalToString(payment.Price)
		p.Value = price
	}

	var pS PaymentStatus
	p.Status = pS.fromDomain(payment.Status)
}

func (pS PaymentStatus) toDomain() domain.PaymentStatus {
	switch pS {
	case PAYMENT_STATUS_OPEN:
		return domain.PAYMENT_STATUS_OPEN
	case PAYMENT_STATUS_APPROVED:
		return domain.PAYMENT_STATUS_APPROVED
	}
	return domain.PAYMENT_SATUS_REFUSED
}

func (pS PaymentStatus) fromDomain(status domain.PaymentStatus) PaymentStatus {
	switch status {
	case domain.PAYMENT_STATUS_OPEN:
		return PAYMENT_STATUS_OPEN
	case domain.PAYMENT_STATUS_APPROVED:
		return PAYMENT_STATUS_APPROVED
	}
	return PAYMENT_STATUS_REFUSED
}

func (pS PaymentStatus) fromRequest(approved bool) PaymentStatus {
	if approved {
		return PAYMENT_STATUS_APPROVED
	}
	return PAYMENT_STATUS_REFUSED
}

func (p *Product) fromDomain(product *domain.Product) {
	if p == nil {
		p = &Product{}
	}

	price := helpers.ParseDecimalToString(product.Price)

	p.ID = product.ID.String()

	if !product.CreatedAt.IsZero() {
		p.CreatedAt = product.CreatedAt.Format(time.RFC3339)
	} else {
		p.CreatedAt = ""
	}
	if !product.UpdatedAt.IsZero() {
		p.UpdatedAt = product.UpdatedAt.Format(time.RFC3339)
	} else {
		p.UpdatedAt = ""
	}
	if !product.DeletedAt.IsZero() {
		p.DeletedAt = product.UpdatedAt.Format(time.RFC3339)
	} else {
		p.DeletedAt = ""
	}

	p.Name = product.Name
	p.Description = product.Description
	p.CategoryID = product.CategoryID.String()
	p.Price = price
}

func (p *Product) toDomain() *domain.Product {
	if p == nil {
		panic("empty product")
	}

	return domain.ParseProductToDomain(
		helpers.SafeUUIDFromString(p.ID),
		helpers.SafeUUIDFromString(p.CategoryID),
		p.Name,
		p.Description,
		p.Price,
	)
}

func (iP *InsertionProduct) toDomain() *domain.Product {
	if iP == nil {
		panic("empty product")
	}

	return domain.NewProduct(uuid.New(), helpers.SafeUUIDFromString(iP.CategoryID), iP.Name, iP.Description, iP.Price)
}

func (uP *UpdateProduct) toDomain() *domain.Product {
	if uP == nil {
		panic("empty product")
	}

	return domain.ParseProductToDomain(
		helpers.SafeUUIDFromString(uP.ID),
		helpers.SafeUUIDFromString(uP.CategoryID),
		uP.Name,
		uP.Description,
		uP.Price,
	)
}

func productsToDomainProducts(products []uuid.UUID) []domain.Product {
	var dPs []domain.Product
	for _, p := range products {
		dP := &domain.Product{ID: p}

		dPs = append(dPs, *dP)
	}
	return dPs
}

func (u *User) fromDomain(user *domain.User) {
	if u == nil {
		u = &User{}
	}

	u.ID = user.ID.String()
	u.Name = user.Name
	u.Document = user.Document
	u.Email = user.Email
}

func (u *User) toDomain() *domain.User {
	if u == nil {
		u = &User{}
	}

	return &domain.User{
		ID:       helpers.SafeUUIDFromString(u.ID),
		Document: u.Document,
		Name:     u.Name,
		Email:    u.Email,
	}
}
