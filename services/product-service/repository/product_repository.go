package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"shopsphere/product-service/models"
)

type ProductRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (r *ProductRepository) Create(
	ctx context.Context,
	product *models.Product,
) (*models.Product, error) {

	query := `
		INSERT INTO products (
			category_id,
			name,
			description,
			price,
			stock_quantity,
			sku,
			image_url,
			is_active
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING
			id,
			category_id,
			name,
			description,
			price,
			stock_quantity,
			sku,
			image_url,
			is_active,
			created_at,
			updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		product.CategoryID,
		product.Name,
		product.Description,
		product.Price,
		product.StockQuantity,
		product.SKU,
		product.ImageURL,
		product.IsActive,
	).Scan(
		&product.ID,
		&product.CategoryID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.StockQuantity,
		&product.SKU,
		&product.ImageURL,
		&product.IsActive,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return product, nil
}

func (r *ProductRepository) FindByID(
	ctx context.Context,
	id int64,
) (*models.Product, error) {

	query := `
		SELECT
			id,
			category_id,
			name,
			description,
			price,
			stock_quantity,
			sku,
			image_url,
			is_active,
			created_at,
			updated_at
		FROM products
		WHERE id = $1
	`

	product := &models.Product{}

	err := r.db.QueryRow(ctx, query, id).Scan(
		&product.ID,
		&product.CategoryID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.StockQuantity,
		&product.SKU,
		&product.ImageURL,
		&product.IsActive,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to find product: %w", err)
	}

	return product, nil
}

func (r *ProductRepository) FindAll(
	ctx context.Context,
) ([]models.Product, error) {

	query := `
		SELECT
			id,
			category_id,
			name,
			description,
			price,
			stock_quantity,
			sku,
			image_url,
			is_active,
			created_at,
			updated_at
		FROM products
		WHERE is_active = TRUE
		ORDER BY id DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch products: %w", err)
	}

	defer rows.Close()

	products := make([]models.Product, 0)

	for rows.Next() {

		var product models.Product

		err := rows.Scan(
			&product.ID,
			&product.CategoryID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.StockQuantity,
			&product.SKU,
			&product.ImageURL,
			&product.IsActive,
			&product.CreatedAt,
			&product.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed while reading products: %w", err)
	}

	return products, nil
}
func (r *ProductRepository) Update(
	ctx context.Context,
	id int64,
	product *models.Product,
) (*models.Product, error) {

	query := `
		UPDATE products
		SET
			category_id = $1,
			name = $2,
			description = $3,
			price = $4,
			stock_quantity = $5,
			sku = $6,
			image_url = $7,
			is_active = $8,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $9
		RETURNING
			id,
			category_id,
			name,
			description,
			price,
			stock_quantity,
			sku,
			image_url,
			is_active,
			created_at,
			updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		product.CategoryID,
		product.Name,
		product.Description,
		product.Price,
		product.StockQuantity,
		product.SKU,
		product.ImageURL,
		product.IsActive,
		id,
	).Scan(
		&product.ID,
		&product.CategoryID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.StockQuantity,
		&product.SKU,
		&product.ImageURL,
		&product.IsActive,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return product, nil
}

func (r *ProductRepository) Delete(
	ctx context.Context,
	id int64,
) (bool, error) {

	query := `
		UPDATE products
		SET
			is_active = FALSE,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING id
	`

	var deletedID int64

	err := r.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(&deletedID)

	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}

	if err != nil {
		return false, fmt.Errorf("failed to delete product: %w", err)
	}

	return true, nil
}
func (r *ProductRepository) FindWithFilters(
	ctx context.Context,
	filter models.ProductFilter,
) ([]models.Product, error) {

	query := `
		SELECT
			id,
			category_id,
			name,
			description,
			price,
			stock_quantity,
			sku,
			image_url,
			is_active,
			created_at,
			updated_at
		FROM products
		WHERE is_active = TRUE
	`

	args := make([]interface{}, 0)
	argIndex := 1

	if strings.TrimSpace(filter.Search) != "" {
		query += fmt.Sprintf(
			" AND (LOWER(name) LIKE LOWER($%d) OR LOWER(description) LIKE LOWER($%d))",
			argIndex,
			argIndex,
		)
		args = append(args, "%"+strings.TrimSpace(filter.Search)+"%")
		argIndex++
	}

	if filter.CategoryID > 0 {
		query += fmt.Sprintf(" AND category_id = $%d", argIndex)
		args = append(args, filter.CategoryID)
		argIndex++
	}

	if filter.MinPrice > 0 {
		query += fmt.Sprintf(" AND price >= $%d", argIndex)
		args = append(args, filter.MinPrice)
		argIndex++
	}

	if filter.MaxPrice > 0 {
		query += fmt.Sprintf(" AND price <= $%d", argIndex)
		args = append(args, filter.MaxPrice)
		argIndex++
	}

	query += " ORDER BY id ASC"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search products: %w", err)
	}

	defer rows.Close()

	products := make([]models.Product, 0)

	for rows.Next() {
		var product models.Product

		err := rows.Scan(
			&product.ID,
			&product.CategoryID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.StockQuantity,
			&product.SKU,
			&product.ImageURL,
			&product.IsActive,
			&product.CreatedAt,
			&product.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed while reading products: %w", err)
	}

	return products, nil
}
