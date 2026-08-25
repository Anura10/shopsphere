package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"shopsphere/product-service/models"
)

type CategoryRepository struct {
	db *pgxpool.Pool
}

func NewCategoryRepository(db *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{
		db: db,
	}
}

func (r *CategoryRepository) Create(
	ctx context.Context,
	category *models.Category,
) (*models.Category, error) {

	query := `
		INSERT INTO categories (
			name,
			description
		)
		VALUES ($1, $2)
		RETURNING
			id,
			name,
			description,
			created_at,
			updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		category.Name,
		category.Description,
	).Scan(
		&category.ID,
		&category.Name,
		&category.Description,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	return category, nil
}

func (r *CategoryRepository) FindByID(
	ctx context.Context,
	id int64,
) (*models.Category, error) {

	query := `
		SELECT
			id,
			name,
			description,
			created_at,
			updated_at
		FROM categories
		WHERE id = $1
	`

	category := &models.Category{}

	err := r.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&category.ID,
		&category.Name,
		&category.Description,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to find category: %w", err)
	}

	return category, nil
}

func (r *CategoryRepository) FindAll(
	ctx context.Context,
) ([]models.Category, error) {

	query := `
		SELECT
			id,
			name,
			description,
			created_at,
			updated_at
		FROM categories
		ORDER BY id ASC
	`

	rows, err := r.db.Query(ctx, query)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch categories: %w", err)
	}

	defer rows.Close()

	categories := make([]models.Category, 0)

	for rows.Next() {

		var category models.Category

		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Description,
			&category.CreatedAt,
			&category.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}

		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed while reading categories: %w", err)
	}

	return categories, nil
}

func (r *CategoryRepository) Update(
	ctx context.Context,
	id int64,
	category *models.Category,
) (*models.Category, error) {

	query := `
		UPDATE categories
		SET
			name = $1,
			description = $2,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
		RETURNING
			id,
			name,
			description,
			created_at,
			updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		category.Name,
		category.Description,
		id,
	).Scan(
		&category.ID,
		&category.Name,
		&category.Description,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	return category, nil
}

func (r *CategoryRepository) Delete(
	ctx context.Context,
	id int64,
) error {

	query := `
		DELETE FROM categories
		WHERE id = $1
	`

	result, err := r.db.Exec(ctx, query, id)

	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	if result.RowsAffected() == 0 {
		return errors.New("category not found")
	}

	return nil
}
