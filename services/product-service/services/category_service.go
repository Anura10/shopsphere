package services

import (
	"context"
	"errors"
	"strings"

	"shopsphere/product-service/models"
	"shopsphere/product-service/repository"
)

type CategoryService struct {
	categoryRepository *repository.CategoryRepository
}

func NewCategoryService(
	categoryRepository *repository.CategoryRepository,
) *CategoryService {
	return &CategoryService{
		categoryRepository: categoryRepository,
	}
}

func (s *CategoryService) Create(
	ctx context.Context,
	req models.CreateCategoryRequest,
) (*models.Category, error) {

	name := strings.TrimSpace(req.Name)

	if name == "" {
		return nil, errors.New("category name is required")
	}

	category := &models.Category{
		Name:        name,
		Description: strings.TrimSpace(req.Description),
	}

	return s.categoryRepository.Create(ctx, category)
}

func (s *CategoryService) GetByID(
	ctx context.Context,
	id int64,
) (*models.Category, error) {

	if id <= 0 {
		return nil, errors.New("invalid category ID")
	}

	return s.categoryRepository.FindByID(ctx, id)
}

func (s *CategoryService) GetAll(
	ctx context.Context,
) ([]models.Category, error) {

	return s.categoryRepository.FindAll(ctx)
}

func (s *CategoryService) Update(
	ctx context.Context,
	id int64,
	req models.UpdateCategoryRequest,
) (*models.Category, error) {

	if id <= 0 {
		return nil, errors.New("invalid category ID")
	}

	existing, err := s.categoryRepository.FindByID(ctx, id)

	if err != nil {
		return nil, err
	}

	if existing == nil {
		return nil, errors.New("category not found")
	}

	if strings.TrimSpace(req.Name) != "" {
		existing.Name = strings.TrimSpace(req.Name)
	}

	existing.Description = strings.TrimSpace(req.Description)

	return s.categoryRepository.Update(
		ctx,
		id,
		existing,
	)
}

func (s *CategoryService) Delete(
	ctx context.Context,
	id int64,
) error {

	if id <= 0 {
		return errors.New("invalid category ID")
	}

	return s.categoryRepository.Delete(ctx, id)
}
