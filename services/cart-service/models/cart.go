package models

import "time"

type Cart struct {
	ID        int64      `json:"id"`
	UserID    int64      `json:"user_id"`
	Items     []CartItem `json:"items"`
	Total     float64    `json:"total"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type CartItem struct {
	ID        int64     `json:"id"`
	CartID    int64     `json:"cart_id"`
	ProductID int64     `json:"product_id"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`

	// Product information returned by the cart API
	ProductName  string  `json:"product_name,omitempty"`
	ProductPrice float64 `json:"product_price,omitempty"`
	Subtotal     float64 `json:"subtotal,omitempty"`
}

type AddCartItemRequest struct {
	ProductID int64 `json:"product_id" binding:"required"`
	Quantity  int   `json:"quantity" binding:"required,min=1"`
}

type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" binding:"required,min=1"`
}