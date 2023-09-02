package http

import (
	"github.com/google/uuid"
	"time"
)

// Multipurpose
type (
	QueryStruct struct {
		ID     string `json:"id" description:"id do objeto a ser removido/consultado"`
		UserID string `json:"user_id" description:"id do usuário requerente"`
	}

	ListRequest struct {
		UserID string `json:"user_id"`
		Limit  int    `json:"limit" default:"10" description:"Quantidade de registros"`
		Offset int    `json:"offset"`
	}
)

// Categories' models
type (
	InsertionCategory struct {
		Name   string `json:"name" description:"Nome da categoria de produto"`
		UserID string `json:"user_id,omitempty" description:"ID do usuario criando categoria"`
	}

	Category struct {
		ID        string    `json:"id" description:"ID da categoria de produto"`
		CreatedAt time.Time `json:"created_at" description:"Epoch time em que categoria foi criada"`
		UpdatedAt time.Time `json:"updated_at" description:"Epoch time em que categoria foi modificada"`
		InsertionCategory
	}

	CategoriesList struct {
		Categories []Category `json:"categories"`
		Limit      int        `json:"limit" default:"10"`
		Offset     int        `json:"offset"`
		Total      int64      `json:"total"`
	}
)

// Orders' models
type (
	Checkout struct {
		Order       Order       `json:"order" description:"Pedido"`
		PaymentInfo PaymentInfo `json:"payment_info" description:"Informações de Cobrança"`
	}

	PaymentInfo struct {
		PaymentID string `json:"payment_id"`
		Value     string `json:"value" description:"Valor a ser pago"`
	}

	OrderStatus string

	Order struct {
		ID        uuid.UUID   `json:"id" description:"ID do Pedido"`
		PaymentID string      `json:"payment_id,omitempty" description:"ID do pagamento"`
		CreatedAt string      `json:"created_at" description:"Data de criação"`
		UpdatedAt string      `json:"updated_at,omitempty" description:"Data de atualização"`
		DeletedAt string      `json:"deleted_at,omitempty" description:"Data de deleção"`
		Price     string      `json:"price" description:"Preço do pedido"`
		Status    OrderStatus `json:"status" description:"Status do pedido"`
		Products  []Product   `json:"products" description:"Lista de Pedidos"`
	}

	InsertionOrder struct {
		UserID      string      `json:"user_id" description:"ID do dono do pedido"`
		ProductsIDs []uuid.UUID `json:"products_ids" description:"ID dos produtos"`
	}

	InsertionOrderSwagger struct {
		UserID      string   `json:"user_id" description:"ID do dono do pedido"`
		ProductsIDs []string `json:"products_ids" description:"Lista de ID dos produtos separados por vírgula"`
	}

	UpdateOrder struct {
		ID string `json:"id" description:"ID do Pedido"`
		InsertionOrder
	}

	OrderList struct {
		Orders []Order `json:"orders"`
		Limit  int     `json:"limit" default:"10"`
		Offset int     `json:"offset"`
		Total  int64   `json:"total"`
	}

	OrderCheckoutRequest struct {
		UserID  string `json:"user_id"`
		OrderID string `json:"order_id" description:"ID do Pedido"`
	}

	OrderStatusUpdate struct {
		OrderID string `json:"order_id" description:"Código de identificação do pedido"`
		UserID  string `json:"user_id" description:"Código de descrição do usuário requerente"`
		Status  string `json:"status" description:"Status para qual deseja mudar o pedido" enum:"Recebido|Preparacao|Pronto|Finalizado|Cancelado"`
	}
)

// Payments' models
type (
	PaymentNotification struct {
		PaymentID string `json:"payment_id" description:"ID do pagamento"`
		OrderID   string `json:"order_id" description:"ID do pedido a ser pago"`
		Approved  bool   `json:"approved" description:"True para aprovado false para recusado"`
	}

	PaymentStatus string

	Payment struct {
		ID        string        `json:"id" description:"ID do pagamento"`
		OrderID   string        `json:"order_id" description:"ID do pedido a ser pago"`
		CreatedAt string        `json:"created_at" description:"Data de criação"`
		UpdatedAt string        `json:"updated_at" description:"Data de Atualização"`
		Value     string        `json:"value" description:"Valor em R$"`
		Status    PaymentStatus `json:"status do pagamento"`
	}
)

// Products' Models
type (
	Product struct {
		UpdateProduct
		InsertionProduct
		CreatedAt string `json:"created_at,omitempty" readOnly:"true"`
		UpdatedAt string `json:"updated_at,omitempty" readOnly:"true"`
		DeletedAt string `json:"deleted_at,omitempty" readOnly:"true"`
	}

	InsertionProduct struct {
		UserID      string `json:"user_id,omitempty"`
		Name        string `json:"name"`
		Description string `json:"description"`
		CategoryID  string `json:"category_id"`
		Price       string `json:"price"`
	}

	UpdateProduct struct {
		InsertionProduct
		ID string `json:"id,omitempty"`
	}

	ProductList struct {
		Products []Product `json:"products"`
		Limit    int       `json:"limit" default:"10"`
		Offset   int       `json:"offset"`
		Total    int64     `json:"total"`
	}
)

//Users' Models
type (
	InsertionUser struct {
		Document string `json:"document" description:"CPF do cliente"`
		Name     string `json:"name" description:"Nome do cliente"`
		Email    string `json:"email" description:"Email do cliente"`
	}

	ValidatedUser struct {
		ID string `json:"id" description:"ID do cliente no sistema"`
	}

	User struct {
		ID string `json:"id" description:"ID do cliente no sistema"`
		InsertionUser
	}

	QueryUser struct {
		Document string `json:"document"`
	}
)

const (
	ORDER_STATUS_UNSET           OrderStatus = ""
	ORDER_STATUS_OPEN                        = "Aberto"
	ORDER_STATUS_WAITING_PAYMENT             = "Aguardando Pagamento"
	ORDER_STATUS_RECEIVED                    = "Recebido"
	ORDER_STATUS_PREPARING                   = "Preparacao"
	ORDER_STATUS_DONE                        = "Pronto"
	ORDER_STATUS_FINISHED                    = "Finalizado"
	ORDER_STATUS_CANCELED                    = "Cancelado"
)

const (
	PAYMENT_STATUS_OPEN     PaymentStatus = "Aguardando Pagamento"
	PAYMENT_STATUS_APPROVED               = "Aprovado"
	PAYMENT_STATUS_REFUSED                = "Recusado"
)
