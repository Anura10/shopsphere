package services

import (
	"context"
	"errors"

	"shopsphere/cart-service/models"
	"shopsphere/cart-service/repository"
)

type CartService struct {
	repository *repository.CartRepository
}

func NewCartService(
	repository *repository.CartRepository,
) *CartService {
	return &CartService{
		repository: repository,
	}
}

func (s *CartService) GetOrCreateCart(
	ctx context.Context,
	userID int64,
) (*models.Cart, error) {

	if userID <= 0 {
		return nil, errors.New("invalid user ID")
	}

	cart, err := s.repository.FindCartByUserID(
		ctx,
		userID,
	)

	if err != nil {
		return nil, err
	}

	if cart == nil {
		cart, err = s.repository.CreateCart(
			ctx,
			userID,
		)

		if err != nil {
			return nil, err
		}
	}

	return s.populateCart(ctx, cart)
}

func (s *CartService) AddItem(
	ctx context.Context,
	userID int64,
	req models.AddCartItemRequest,
) (*models.Cart, error) {

	if userID <= 0 {
		return nil, errors.New("invalid user ID")
	}

	if req.ProductID <= 0 {
		return nil, errors.New("invalid product ID")
	}

	if req.Quantity <= 0 {
		return nil, errors.New("quantity must be greater than zero")
	}

	cart, err := s.GetOrCreateCart(
		ctx,
		userID,
	)

	if err != nil {
		return nil, err
	}

	existing, err := s.repository.FindCartItem(
		ctx,
		cart.ID,
		req.ProductID,
	)

	if err != nil {
		return nil, err
	}

	if existing != nil {

		_, err = s.repository.UpdateItem(
			ctx,
			cart.ID,
			req.ProductID,
			existing.Quantity+req.Quantity,
		)

		if err != nil {
			return nil, err
		}

	} else {

		_, err = s.repository.AddItem(
			ctx,
			cart.ID,
			req.ProductID,
			req.Quantity,
		)

		if err != nil {
			return nil, err
		}
	}

	return s.populateCart(ctx, cart)
}

func (s *CartService) UpdateItem(
	ctx context.Context,
	userID int64,
	productID int64,
	req models.UpdateCartItemRequest,
) (*models.Cart, error) {

	if userID <= 0 {
		return nil, errors.New("invalid user ID")
	}

	if productID <= 0 {
		return nil, errors.New("invalid product ID")
	}

	if req.Quantity <= 0 {
		return nil, errors.New("quantity must be greater than zero")
	}

	cart, err := s.GetOrCreateCart(
		ctx,
		userID,
	)

	if err != nil {
		return nil, err
	}

	item, err := s.repository.FindCartItem(
		ctx,
		cart.ID,
		productID,
	)

	if err != nil {
		return nil, err
	}

	if item == nil {
		return nil, errors.New("cart item not found")
	}

	_, err = s.repository.UpdateItem(
		ctx,
		cart.ID,
		productID,
		req.Quantity,
	)

	if err != nil {
		return nil, err
	}

	return s.populateCart(ctx, cart)
}

func (s *CartService) RemoveItem(
	ctx context.Context,
	userID int64,
	productID int64,
) (*models.Cart, error) {

	cart, err := s.GetOrCreateCart(
		ctx,
		userID,
	)

	if err != nil {
		return nil, err
	}

	err = s.repository.DeleteItem(
		ctx,
		cart.ID,
		productID,
	)

	if err != nil {
		return nil, err
	}

	return s.populateCart(ctx, cart)
}

func (s *CartService) ClearCart(
	ctx context.Context,
	userID int64,
) (*models.Cart, error) {

	cart, err := s.GetOrCreateCart(
		ctx,
		userID,
	)

	if err != nil {
		return nil, err
	}

	err = s.repository.ClearCart(
		ctx,
		cart.ID,
	)

	if err != nil {
		return nil, err
	}

	return s.populateCart(ctx, cart)
}

func (s *CartService) populateCart(
	ctx context.Context,
	cart *models.Cart,
) (*models.Cart, error) {

	items, err := s.repository.GetCartItems(
		ctx,
		cart.ID,
	)

	if err != nil {
		return nil, err
	}

	cart.Items = items
	cart.Total = 0

	for _, item := range items {
		cart.Total += item.Subtotal
	}

	return cart, nil
}