package models

import "time"

type Inventory struct {
	ID               int64     `json:"id"`
	ProductID        int64     `json:"product_id"`
	Quantity         int       `json:"quantity"`
	ReservedQuantity int       `json:"reserved_quantity"`
	ReorderLevel     int       `json:"reorder_level"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type CreateInventoryRequest struct {
	ProductID    int64 `json:"product_id" binding:"required"`
	Quantity     int   `json:"quantity"`
	ReorderLevel int   `json:"reorder_level"`
}

type UpdateInventoryRequest struct {
	Quantity     int `json:"quantity"`
	ReorderLevel int `json:"reorder_level"`
}

type StockRequest struct {
	Quantity int `json:"quantity" binding:"required,min=1"`
}
