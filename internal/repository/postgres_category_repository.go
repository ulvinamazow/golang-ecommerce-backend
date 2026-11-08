package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ulvinamazow/go-backend-project/internal/model"
)

type postgresCategoryRepository struct {
	db *pgxpool.Pool
}

func NewPostgresCategoryRepository(db *pgxpool.Pool) CategoryRepository {
	return &postgresCategoryRepository{db: db}
}

func (r *postgresCategoryRepository) Create(ctx context.Context, category *model.Category) error {
	query := `
		INSERT INTO categories (name, created_at, updated_at)
		VALUES ($1, NOW(), NOW())
		RETURNING name, created_at, updated_at
	`
	err := r.db.QueryRow(ctx, query, category.Name).
		Scan(&category.Name, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		return fmt.Errorf("Category could not be created: %w", err)
	}
	return nil
}

func (r *postgresCategoryRepository) GetByID(ctx context.Context, id int64) (*model.Category, error) {
	query := `
		SELECT id, name, created_at, updated_at
		FROM categories
		WHERE id = $1
	`
	var c model.Category
	err := r.db.QueryRow(ctx, query, id).
		Scan(&c.ID, &c.Name, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("category could not be retrieved: %w", err)
	}
	return &c, nil
}

func (r *postgresCategoryRepository) Update(ctx context.Context, category *model.Category) error {
	query := `
		UPDATE categories
		SET name = $1, updated_at = NOW()
		WHERE id = $2
		RETURNING updated_at
	`
	err := r.db.QueryRow(ctx, query, category.Name, category.ID).
		Scan(&category.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return fmt.Errorf("category could not be updated: %w", err)
	}
	return nil
}

func (r *postgresCategoryRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM categories WHERE id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("category could not be deleted: %w", err)
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *postgresCategoryRepository) List(ctx context.Context, limit, offset int) ([]*model.Category, error) {
	query := `
		SELECT id, name, created_at, updated_at
		FROM categories
		ORDER BY id
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list of categories could not be retrieved: %w", err)
	}
	defer rows.Close()

	categories := make([]*model.Category, 0)
	for rows.Next() {
		var c model.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, fmt.Errorf("error while reading category sequence: %w", err)
		}
		categories = append(categories, &c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while processing category ranges: %w", err)
	}
	return categories, nil
}
