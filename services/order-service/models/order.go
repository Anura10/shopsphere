package models

import "time"

// Order represents a customer order.
type Order struct {
	ID              int64       `json:"id"`
	UserID          int64       `json:"user_id"`
	Status          string      `json:"status"`
	TotalAmount     float64     `json:"total_amount"`
	ShippingAddress string      `json:"shipping_address"`
	Items           []OrderItem `json:"items"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}

// OrderItem represents a product inside an order.
type OrderItem struct {
	ID          int64   `json:"id"`
	OrderID     int64   `json:"order_id"`
	ProductID   int64   `json:"product_id"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	Subtotal    float64 `json:"subtotal"`
	ProductName string  `json:"product_name,omitempty"`
}

// Cart represents the response returned by Cart Service.
type Cart struct {
	ID        int64      `json:"id"`
	UserID    int64      `json:"user_id"`
	Items     []CartItem `json:"items"`
	Total     float64    `json:"total"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// CartItem represents an item returned by Cart Service.
type CartItem struct {
	ID           int64     `json:"id"`
	CartID       int64     `json:"cart_id"`
	ProductID    int64     `json:"product_id"`
	Quantity     int       `json:"quantity"`
	CreatedAt    time.Time `json:"created_at"`
	ProductName  string    `json:"product_name,omitempty"`
	ProductPrice float64   `json:"product_price,omitempty"`
	Subtotal     float64   `json:"subtotal,omitempty"`
}

// CreateOrderRequest contains the information required to create an order.
type CreateOrderRequest struct {
	ShippingAddress string `json:"shipping_address" binding:"required"`
}

const (
	OrderStatusPending   = "PENDING"
	OrderStatusConfirmed = "CONFIRMED"
	OrderStatusCancelled = "CANCELLED"
)
