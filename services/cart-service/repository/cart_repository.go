package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"shopsphere/cart-service/models"
)

type CartRepository struct {
	db *pgxpool.Pool
}

func NewCartRepository(db *pgxpool.Pool) *CartRepository {
	return &CartRepository{
		db: db,
	}
}

func (r *CartRepository) CreateCart(
	ctx context.Context,
	userID int64,
) (*models.Cart, error) {

	query := `
		INSERT INTO carts (user_id)
		VALUES ($1)
		RETURNING id, user_id, created_at, updated_at
	`

	cart := &models.Cart{}

	err := r.db.QueryRow(
		ctx,
		query,
		userID,
	).Scan(
		&cart.ID,
		&cart.UserID,
		&cart.CreatedAt,
		&cart.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create cart: %w", err)
	}

	cart.Items = []models.CartItem{}

	return cart, nil
}

func (r *CartRepository) FindCartByUserID(
	ctx context.Context,
	userID int64,
) (*models.Cart, error) {

	query := `
		SELECT
			id,
			user_id,
			created_at,
			updated_at
		FROM carts
		WHERE user_id = $1
	`

	cart := &models.Cart{}

	err := r.db.QueryRow(
		ctx,
		query,
		userID,
	).Scan(
		&cart.ID,
		&cart.UserID,
		&cart.CreatedAt,
		&cart.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to find cart: %w", err)
	}

	return cart, nil
}

func (r *CartRepository) GetCartItems(
	ctx context.Context,
	cartID int64,
) ([]models.CartItem, error) {

	query := `
		SELECT
			ci.id,
			ci.cart_id,
			ci.product_id,
			ci.quantity,
			ci.created_at,
			p.name,
			p.price
		FROM cart_items ci
		JOIN products p
			ON p.id = ci.product_id
		WHERE ci.cart_id = $1
		ORDER BY ci.id ASC
	`

	rows, err := r.db.Query(
		ctx,
		query,
		cartID,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get cart items: %w", err)
	}

	defer rows.Close()

	items := make([]models.CartItem, 0)

	for rows.Next() {

		var item models.CartItem

		err := rows.Scan(
			&item.ID,
			&item.CartID,
			&item.ProductID,
			&item.Quantity,
			&item.CreatedAt,
			&item.ProductName,
			&item.ProductPrice,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan cart item: %w", err)
		}

		item.Subtotal = item.ProductPrice * float64(item.Quantity)

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed while reading cart items: %w", err)
	}

	return items, nil
}

func (r *CartRepository) FindCartItem(
	ctx context.Context,
	cartID int64,
	productID int64,
) (*models.CartItem, error) {

	query := `
		SELECT
			id,
			cart_id,
			product_id,
			quantity,
			created_at
		FROM cart_items
		WHERE cart_id = $1
		  AND product_id = $2
	`

	item := &models.CartItem{}

	err := r.db.QueryRow(
		ctx,
		query,
		cartID,
		productID,
	).Scan(
		&item.ID,
		&item.CartID,
		&item.ProductID,
		&item.Quantity,
		&item.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to find cart item: %w", err)
	}

	return item, nil
}

func (r *CartRepository) AddItem(
	ctx context.Context,
	cartID int64,
	productID int64,
	quantity int,
) (*models.CartItem, error) {

	query := `
		INSERT INTO cart_items (
			cart_id,
			product_id,
			quantity
		)
		VALUES ($1, $2, $3)
		RETURNING
			id,
			cart_id,
			product_id,
			quantity,
			created_at
	`

	item := &models.CartItem{}

	err := r.db.QueryRow(
		ctx,
		query,
		cartID,
		productID,
		quantity,
	).Scan(
		&item.ID,
		&item.CartID,
		&item.ProductID,
		&item.Quantity,
		&item.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to add cart item: %w", err)
	}

	return item, nil
}

func (r *CartRepository) UpdateItem(
	ctx context.Context,
	cartID int64,
	productID int64,
	quantity int,
) (*models.CartItem, error) {

	query := `
		UPDATE cart_items
		SET quantity = $1
		WHERE cart_id = $2
		  AND product_id = $3
		RETURNING
			id,
			cart_id,
			product_id,
			quantity,
			created_at
	`

	item := &models.CartItem{}

	err := r.db.QueryRow(
		ctx,
		query,
		quantity,
		cartID,
		productID,
	).Scan(
		&item.ID,
		&item.CartID,
		&item.ProductID,
		&item.Quantity,
		&item.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to update cart item: %w", err)
	}

	return item, nil
}

func (r *CartRepository) DeleteItem(
	ctx context.Context,
	cartID int64,
	productID int64,
) error {

	query := `
		DELETE FROM cart_items
		WHERE cart_id = $1
		  AND product_id = $2
	`

	result, err := r.db.Exec(
		ctx,
		query,
		cartID,
		productID,
	)

	if err != nil {
		return fmt.Errorf("failed to delete cart item: %w", err)
	}

	if result.RowsAffected() == 0 {
		return errors.New("cart item not found")
	}

	return nil
}

func (r *CartRepository) ClearCart(
	ctx context.Context,
	cartID int64,
) error {

	query := `
		DELETE FROM cart_items
		WHERE cart_id = $1
	`

	_, err := r.db.Exec(
		ctx,
		query,
		cartID,
	)

	if err != nil {
		return fmt.Errorf("failed to clear cart: %w", err)
	}

	return nil
}