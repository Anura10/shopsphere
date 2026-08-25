package models

import "time"

type Product struct {
	ID            int64     `json:"id"`
	CategoryID    int64     `json:"category_id"`
	Name          string    `json:"name"`
	Description   string    `json:"description,omitempty"`
	Price         float64   `json:"price"`
	StockQuantity int       `json:"stock_quantity"`
	SKU           string    `json:"sku"`
	ImageURL      string    `json:"image_url,omitempty"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type CreateProductRequest struct {
	CategoryID    int64   `json:"category_id" binding:"required"`
	Name          string  `json:"name" binding:"required"`
	Description   string  `json:"description"`
	Price         float64 `json:"price" binding:"required"`
	StockQuantity int     `json:"stock_quantity"`
	SKU           string  `json:"sku" binding:"required"`
	ImageURL      string  `json:"image_url"`
}

type UpdateProductRequest struct {
	CategoryID    int64   `json:"category_id"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Price         float64 `json:"price"`
	StockQuantity int     `json:"stock_quantity"`
	SKU           string  `json:"sku"`
	ImageURL      string  `json:"image_url"`
	IsActive      *bool   `json:"is_active"`
}
type ProductFilter struct {
	Search     string
	CategoryID int64
	MinPrice   float64
	MaxPrice   float64
}
