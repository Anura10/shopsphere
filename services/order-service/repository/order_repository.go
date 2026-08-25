package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"shopsphere/order-service/models"
)

type OrderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (r *OrderRepository) CreateOrder(
	ctx context.Context,
	userID int64,
	shippingAddress string,
	totalAmount float64,
) (*models.Order, error) {

	query := `
		INSERT INTO orders (
			user_id,
			status,
			total_amount,
			shipping_address
		)
		VALUES ($1, $2, $3, $4)
		RETURNING
			id,
			user_id,
			status,
			total_amount,
			shipping_address,
			created_at,
			updated_at
	`

	order := &models.Order{}

	err := r.db.QueryRow(
		ctx,
		query,
		userID,
		models.OrderStatusPending,
		totalAmount,
		shippingAddress,
	).Scan(
		&order.ID,
		&order.UserID,
		&order.Status,
		&order.TotalAmount,
		&order.ShippingAddress,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to create order: %w",
			err,
		)
	}

	order.Items = []models.OrderItem{}

	return order, nil
}

func (r *OrderRepository) AddOrderItem(
	ctx context.Context,
	orderID int64,
	productID int64,
	quantity int,
	unitPrice float64,
	productName string,
) (*models.OrderItem, error) {

	query := `
		INSERT INTO order_items (
			order_id,
			product_id,
			quantity,
			unit_price,
			subtotal,
			product_name
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING
			id,
			order_id,
			product_id,
			quantity,
			unit_price,
			subtotal,
			product_name
	`

	subtotal := unitPrice * float64(quantity)

	item := &models.OrderItem{}

	err := r.db.QueryRow(
		ctx,
		query,
		orderID,
		productID,
		quantity,
		unitPrice,
		subtotal,
		productName,
	).Scan(
		&item.ID,
		&item.OrderID,
		&item.ProductID,
		&item.Quantity,
		&item.UnitPrice,
		&item.Subtotal,
		&item.ProductName,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to add order item: %w",
			err,
		)
	}

	return item, nil
}

func (r *OrderRepository) GetOrderByID(
	ctx context.Context,
	orderID int64,
	userID int64,
) (*models.Order, error) {

	query := `
		SELECT
			id,
			user_id,
			status,
			total_amount,
			shipping_address,
			created_at,
			updated_at
		FROM orders
		WHERE id = $1
		  AND user_id = $2
	`

	order := &models.Order{}

	err := r.db.QueryRow(
		ctx,
		query,
		orderID,
		userID,
	).Scan(
		&order.ID,
		&order.UserID,
		&order.Status,
		&order.TotalAmount,
		&order.ShippingAddress,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf(
			"failed to get order: %w",
			err,
		)
	}

	items, err := r.GetOrderItems(
		ctx,
		order.ID,
	)

	if err != nil {
		return nil, err
	}

	order.Items = items

	return order, nil
}

func (r *OrderRepository) GetOrderItems(
	ctx context.Context,
	orderID int64,
) ([]models.OrderItem, error) {

	query := `
		SELECT
			id,
			order_id,
			product_id,
			quantity,
			unit_price,
			subtotal,
			product_name
		FROM order_items
		WHERE order_id = $1
		ORDER BY id ASC
	`

	rows, err := r.db.Query(
		ctx,
		query,
		orderID,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to get order items: %w",
			err,
		)
	}

	defer rows.Close()

	items := make([]models.OrderItem, 0)

	for rows.Next() {

		var item models.OrderItem

		err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.ProductID,
			&item.Quantity,
			&item.UnitPrice,
			&item.Subtotal,
			&item.ProductName,
		)

		if err != nil {
			return nil, fmt.Errorf(
				"failed to scan order item: %w",
				err,
			)
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"failed while reading order items: %w",
			err,
		)
	}

	return items, nil
}

func (r *OrderRepository) GetOrdersByUserID(
	ctx context.Context,
	userID int64,
) ([]models.Order, error) {

	query := `
		SELECT
			id,
			user_id,
			status,
			total_amount,
			shipping_address,
			created_at,
			updated_at
		FROM orders
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(
		ctx,
		query,
		userID,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to get orders: %w",
			err,
		)
	}

	defer rows.Close()

	orders := make([]models.Order, 0)

	for rows.Next() {

		var order models.Order

		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.Status,
			&order.TotalAmount,
			&order.ShippingAddress,
			&order.CreatedAt,
			&order.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf(
				"failed to scan order: %w",
				err,
			)
		}

		items, err := r.GetOrderItems(
			ctx,
			order.ID,
		)

		if err != nil {
			return nil, err
		}

		order.Items = items

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"failed while reading orders: %w",
			err,
		)
	}

	return orders, nil
}

func (r *OrderRepository) CancelOrder(
	ctx context.Context,
	orderID int64,
	userID int64,
) (*models.Order, error) {

	query := `
		UPDATE orders
		SET
			status = $1,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
		  AND user_id = $3
		  AND status = $4
		RETURNING
			id,
			user_id,
			status,
			total_amount,
			shipping_address,
			created_at,
			updated_at
	`

	order := &models.Order{}

	err := r.db.QueryRow(
		ctx,
		query,
		models.OrderStatusCancelled,
		orderID,
		userID,
		models.OrderStatusPending,
	).Scan(
		&order.ID,
		&order.UserID,
		&order.Status,
		&order.TotalAmount,
		&order.ShippingAddress,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf(
			"failed to cancel order: %w",
			err,
		)
	}

	items, err := r.GetOrderItems(
		ctx,
		order.ID,
	)

	if err != nil {
		return nil, err
	}

	order.Items = items

	return order, nil
}
