package services

import (
	"context"
	"errors"
	"strings"

	"shopsphere/product-service/models"
	"shopsphere/product-service/repository"
)

type ProductService struct {
	productRepository *repository.ProductRepository
}

func NewProductService(
	productRepository *repository.ProductRepository,
) *ProductService {
	return &ProductService{
		productRepository: productRepository,
	}
}

func (s *ProductService) Create(
	ctx context.Context,
	req models.CreateProductRequest,
) (*models.Product, error) {

	if strings.TrimSpace(req.Name) == "" {
		return nil, errors.New("product name is required")
	}

	if strings.TrimSpace(req.SKU) == "" {
		return nil, errors.New("SKU is required")
	}

	if req.Price < 0 {
		return nil, errors.New("price cannot be negative")
	}

	if req.StockQuantity < 0 {
		return nil, errors.New("stock quantity cannot be negative")
	}

	product := &models.Product{
		CategoryID:    req.CategoryID,
		Name:          strings.TrimSpace(req.Name),
		Description:   strings.TrimSpace(req.Description),
		Price:         req.Price,
		StockQuantity: req.StockQuantity,
		SKU:           strings.TrimSpace(req.SKU),
		ImageURL:      strings.TrimSpace(req.ImageURL),
		IsActive:      true,
	}

	return s.productRepository.Create(ctx, product)
}

func (s *ProductService) GetByID(
	ctx context.Context,
	id int64,
) (*models.Product, error) {

	if id <= 0 {
		return nil, errors.New("invalid product ID")
	}

	return s.productRepository.FindByID(ctx, id)
}

func (s *ProductService) GetAll(
	ctx context.Context,
) ([]models.Product, error) {

	return s.productRepository.FindAll(ctx)
}
func (s *ProductService) Update(
	ctx context.Context,
	id int64,
	req models.UpdateProductRequest,
) (*models.Product, error) {

	if id <= 0 {
		return nil, errors.New("invalid product ID")
	}

	existing, err := s.productRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if existing == nil {
		return nil, errors.New("product not found")
	}

	if strings.TrimSpace(req.Name) != "" {
		existing.Name = strings.TrimSpace(req.Name)
	}

	if req.CategoryID > 0 {
		existing.CategoryID = req.CategoryID
	}

	if req.Description != "" {
		existing.Description = strings.TrimSpace(req.Description)
	}

	if req.Price >= 0 {
		existing.Price = req.Price
	}

	if req.StockQuantity >= 0 {
		existing.StockQuantity = req.StockQuantity
	}

	if strings.TrimSpace(req.SKU) != "" {
		existing.SKU = strings.TrimSpace(req.SKU)
	}

	if req.ImageURL != "" {
		existing.ImageURL = strings.TrimSpace(req.ImageURL)
	}

	if req.IsActive != nil {
		existing.IsActive = *req.IsActive
	}

	return s.productRepository.Update(ctx, id, existing)
}

func (s *ProductService) Delete(
	ctx context.Context,
	id int64,
) error {

	if id <= 0 {
		return errors.New("invalid product ID")
	}

	deleted, err := s.productRepository.Delete(ctx, id)
	if err != nil {
		return err
	}

	if !deleted {
		return errors.New("product not found")
	}

	return nil
}
func (s *ProductService) Search(
	ctx context.Context,
	filter models.ProductFilter,
) ([]models.Product, error) {

	return s.productRepository.FindWithFilters(ctx, filter)
}
