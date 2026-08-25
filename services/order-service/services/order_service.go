package services

import (
	"context"
	"errors"

	"shopsphere/order-service/models"
	"shopsphere/order-service/repository"
)

type OrderService struct {
	repository      *repository.OrderRepository
	cartClient      *CartClient
	inventoryClient *InventoryClient
}

func NewOrderService(
	repository *repository.OrderRepository,
	cartClient *CartClient,
	inventoryClient *InventoryClient,
) *OrderService {
	return &OrderService{
		repository:      repository,
		cartClient:      cartClient,
		inventoryClient: inventoryClient,
	}
}

func (s *OrderService) CreateOrder(
	ctx context.Context,
	userID int64,
	token string,
	req models.CreateOrderRequest,
) (*models.Order, error) {

	if userID <= 0 {
		return nil, errors.New("invalid user ID")
	}

	if token == "" {
		return nil, errors.New("authentication token is required")
	}

	if req.ShippingAddress == "" {
		return nil, errors.New("shipping address is required")
	}

	// Get the current user's cart.
	cart, err := s.cartClient.GetCart(
		ctx,
		token,
	)
	if err != nil {
		return nil, err
	}

	if cart == nil {
		return nil, errors.New("cart not found")
	}

	if len(cart.Items) == 0 {
		return nil, errors.New("cart is empty")
	}

	// Re-check inventory for EVERY cart item before
	// creating the order.
	for _, item := range cart.Items {

		available, err := s.inventoryClient.HasEnoughStock(
			ctx,
			item.ProductID,
			item.Quantity,
		)

		if err != nil {
			return nil, err
		}

		if !available {
			return nil, errors.New(
				"insufficient stock for product",
			)
		}
	}

	// Only create the order after all inventory checks pass.
	order, err := s.repository.CreateOrder(
		ctx,
		userID,
		req.ShippingAddress,
		cart.Total,
	)

	if err != nil {
		return nil, err
	}

	// Copy cart items into order_items.
	for _, cartItem := range cart.Items {

		_, err := s.repository.AddOrderItem(
			ctx,
			order.ID,
			cartItem.ProductID,
			cartItem.Quantity,
			cartItem.ProductPrice,
			cartItem.ProductName,
		)

		if err != nil {
			return nil, err
		}
	}

	// Load the complete order.
	order, err = s.repository.GetOrderByID(
		ctx,
		order.ID,
		userID,
	)

	if err != nil {
		return nil, err
	}

	return order, nil
}

func (s *OrderService) GetOrder(
	ctx context.Context,
	userID int64,
	orderID int64,
) (*models.Order, error) {

	if userID <= 0 {
		return nil, errors.New("invalid user ID")
	}

	if orderID <= 0 {
		return nil, errors.New("invalid order ID")
	}

	order, err := s.repository.GetOrderByID(
		ctx,
		orderID,
		userID,
	)

	if err != nil {
		return nil, err
	}

	if order == nil {
		return nil, errors.New("order not found")
	}

	return order, nil
}

func (s *OrderService) GetUserOrders(
	ctx context.Context,
	userID int64,
) ([]models.Order, error) {

	if userID <= 0 {
		return nil, errors.New("invalid user ID")
	}

	return s.repository.GetOrdersByUserID(
		ctx,
		userID,
	)
}

func (s *OrderService) CancelOrder(
	ctx context.Context,
	userID int64,
	orderID int64,
) (*models.Order, error) {

	if userID <= 0 {
		return nil, errors.New("invalid user ID")
	}

	if orderID <= 0 {
		return nil, errors.New("invalid order ID")
	}

	order, err := s.repository.CancelOrder(
		ctx,
		orderID,
		userID,
	)

	if err != nil {
		return nil, err
	}

	if order == nil {
		return nil, errors.New(
			"order not found or cannot be cancelled",
		)
	}

	return order, nil
}
