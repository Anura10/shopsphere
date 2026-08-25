package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"shopsphere/inventory-service/models"
)

type InventoryRepository struct {
	db *pgxpool.Pool
}

func NewInventoryRepository(db *pgxpool.Pool) *InventoryRepository {
	return &InventoryRepository{
		db: db,
	}
}

func (r *InventoryRepository) Create(
	ctx context.Context,
	inventory *models.Inventory,
) (*models.Inventory, error) {

	query := `
		INSERT INTO inventory (
			product_id,
			quantity,
			reorder_level
		)
		VALUES ($1, $2, $3)
		RETURNING
			id,
			product_id,
			quantity,
			reserved_quantity,
			reorder_level,
			created_at,
			updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		inventory.ProductID,
		inventory.Quantity,
		inventory.ReorderLevel,
	).Scan(
		&inventory.ID,
		&inventory.ProductID,
		&inventory.Quantity,
		&inventory.ReservedQuantity,
		&inventory.ReorderLevel,
		&inventory.CreatedAt,
		&inventory.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create inventory: %w", err)
	}

	return inventory, nil
}

func (r *InventoryRepository) FindByProductID(
	ctx context.Context,
	productID int64,
) (*models.Inventory, error) {

	query := `
		SELECT
			id,
			product_id,
			quantity,
			reserved_quantity,
			reorder_level,
			created_at,
			updated_at
		FROM inventory
		WHERE product_id = $1
	`

	inventory := &models.Inventory{}

	err := r.db.QueryRow(
		ctx,
		query,
		productID,
	).Scan(
		&inventory.ID,
		&inventory.ProductID,
		&inventory.Quantity,
		&inventory.ReservedQuantity,
		&inventory.ReorderLevel,
		&inventory.CreatedAt,
		&inventory.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to find inventory: %w", err)
	}

	return inventory, nil
}

func (r *InventoryRepository) Update(
	ctx context.Context,
	productID int64,
	quantity int,
	reorderLevel int,
) (*models.Inventory, error) {

	query := `
		UPDATE inventory
		SET
			quantity = $1,
			reorder_level = $2,
			updated_at = CURRENT_TIMESTAMP
		WHERE product_id = $3
		RETURNING
			id,
			product_id,
			quantity,
			reserved_quantity,
			reorder_level,
			created_at,
			updated_at
	`

	inventory := &models.Inventory{}

	err := r.db.QueryRow(
		ctx,
		query,
		quantity,
		reorderLevel,
		productID,
	).Scan(
		&inventory.ID,
		&inventory.ProductID,
		&inventory.Quantity,
		&inventory.ReservedQuantity,
		&inventory.ReorderLevel,
		&inventory.CreatedAt,
		&inventory.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to update inventory: %w", err)
	}

	return inventory, nil
}

func (r *InventoryRepository) AddStock(
	ctx context.Context,
	productID int64,
	amount int,
) (*models.Inventory, error) {

	query := `
		UPDATE inventory
		SET
			quantity = quantity + $1,
			updated_at = CURRENT_TIMESTAMP
		WHERE product_id = $2
		RETURNING
			id,
			product_id,
			quantity,
			reserved_quantity,
			reorder_level,
			created_at,
			updated_at
	`

	inventory := &models.Inventory{}

	err := r.db.QueryRow(
		ctx,
		query,
		amount,
		productID,
	).Scan(
		&inventory.ID,
		&inventory.ProductID,
		&inventory.Quantity,
		&inventory.ReservedQuantity,
		&inventory.ReorderLevel,
		&inventory.CreatedAt,
		&inventory.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to add stock: %w", err)
	}

	return inventory, nil
}

func (r *InventoryRepository) RemoveStock(
	ctx context.Context,
	productID int64,
	amount int,
) (*models.Inventory, error) {

	query := `
		UPDATE inventory
		SET
			quantity = quantity - $1,
			updated_at = CURRENT_TIMESTAMP
		WHERE product_id = $2
		  AND quantity >= $1
		RETURNING
			id,
			product_id,
			quantity,
			reserved_quantity,
			reorder_level,
			created_at,
			updated_at
	`

	inventory := &models.Inventory{}

	err := r.db.QueryRow(
		ctx,
		query,
		amount,
		productID,
	).Scan(
		&inventory.ID,
		&inventory.ProductID,
		&inventory.Quantity,
		&inventory.ReservedQuantity,
		&inventory.ReorderLevel,
		&inventory.CreatedAt,
		&inventory.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to remove stock: %w", err)
	}

	return inventory, nil
}
