package services

import (
	"context"
	"errors"

	"shopsphere/inventory-service/models"
	"shopsphere/inventory-service/repository"
)

type InventoryService struct {
	repository *repository.InventoryRepository
}

func NewInventoryService(
	repository *repository.InventoryRepository,
) *InventoryService {
	return &InventoryService{
		repository: repository,
	}
}

func (s *InventoryService) Create(
	ctx context.Context,
	req models.CreateInventoryRequest,
) (*models.Inventory, error) {

	if req.ProductID <= 0 {
		return nil, errors.New("invalid product ID")
	}

	if req.Quantity < 0 {
		return nil, errors.New("quantity cannot be negative")
	}

	if req.ReorderLevel < 0 {
		return nil, errors.New("reorder level cannot be negative")
	}

	existing, err := s.repository.FindByProductID(
		ctx,
		req.ProductID,
	)

	if err != nil {
		return nil, err
	}

	if existing != nil {
		return nil, errors.New("inventory already exists for this product")
	}

	inventory := &models.Inventory{
		ProductID:    req.ProductID,
		Quantity:     req.Quantity,
		ReorderLevel: req.ReorderLevel,
	}

	return s.repository.Create(ctx, inventory)
}

func (s *InventoryService) GetByProductID(
	ctx context.Context,
	productID int64,
) (*models.Inventory, error) {

	if productID <= 0 {
		return nil, errors.New("invalid product ID")
	}

	return s.repository.FindByProductID(ctx, productID)
}

func (s *InventoryService) Update(
	ctx context.Context,
	productID int64,
	req models.UpdateInventoryRequest,
) (*models.Inventory, error) {

	if productID <= 0 {
		return nil, errors.New("invalid product ID")
	}

	if req.Quantity < 0 {
		return nil, errors.New("quantity cannot be negative")
	}

	if req.ReorderLevel < 0 {
		return nil, errors.New("reorder level cannot be negative")
	}

	return s.repository.Update(
		ctx,
		productID,
		req.Quantity,
		req.ReorderLevel,
	)
}

func (s *InventoryService) AddStock(
	ctx context.Context,
	productID int64,
	amount int,
) (*models.Inventory, error) {

	if productID <= 0 {
		return nil, errors.New("invalid product ID")
	}

	if amount <= 0 {
		return nil, errors.New("quantity must be greater than zero")
	}

	return s.repository.AddStock(
		ctx,
		productID,
		amount,
	)
}

func (s *InventoryService) RemoveStock(
	ctx context.Context,
	productID int64,
	amount int,
) (*models.Inventory, error) {

	if productID <= 0 {
		return nil, errors.New("invalid product ID")
	}

	if amount <= 0 {
		return nil, errors.New("quantity must be greater than zero")
	}

	inventory, err := s.repository.RemoveStock(
		ctx,
		productID,
		amount,
	)

	if err != nil {
		return nil, err
	}

	if inventory == nil {
		return nil, errors.New("insufficient stock or inventory not found")
	}

	return inventory, nil
}
