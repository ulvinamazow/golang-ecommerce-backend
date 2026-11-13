package productrepository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ulvinamazow/go-backend-project/internal/model"
	"github.com/ulvinamazow/go-backend-project/internal/repository"
)

type postgresProductRepository struct {
	db *pgxpool.Pool
}

func NewPostgresProductRepository(db *pgxpool.Pool) ProductRepository {
	return &postgresProductRepository{db: db}
}

func (r *postgresProductRepository) Create(ctx context.Context, product *model.Product) error {
	query := `
		INSERT INTO products (name, description, price, category_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRow(ctx, query, product.Name, product.Description, product.Price, product.CategoryID).
		Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		return fmt.Errorf("product can not be created: %w", err)
	}
	return nil
}

func (r *postgresProductRepository) GetByID(ctx context.Context, id int64) (*model.Product, error) {
	query := `
		SELECT id, name, description, price, category_id, created_at, updated_at
		FROM products
		WHERE id = $1
	`
	var p model.Product
	err := r.db.QueryRow(ctx, query, id).
		Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.CategoryID, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("product could not be brought: %w", err)
	}
	return &p, nil
}

func (r *postgresProductRepository) Update(ctx context.Context, product *model.Product) error {
	query := `
		UPDATE products
		SET name = $1, description = $2, price = $3, category_id = $4, updated_at = NOW()
		WHERE id = $5
		RETURNING updated_at
	`
	err := r.db.QueryRow(ctx, query, product.Name, product.Description, product.Price, product.CategoryID, product.ID).
		Scan(&product.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repository.ErrNotFound
		}
		return fmt.Errorf("product can not be updated: %w", err)
	}
	return nil
}

func (r *postgresProductRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM products WHERE id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("product can not be deleted: %w", err)
	}
	if result.RowsAffected() == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *postgresProductRepository) List(ctx context.Context, limit, offset int) ([]*model.Product, error) {
	query := `
		SELECT id, name, description, price, category_id, created_at, updated_at
		FROM products
		ORDER BY id
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("the list of products could not be retrieved: %w", err)
	}
	defer rows.Close()

	products := make([]*model.Product, 0)
	for rows.Next() {
		var p model.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.CategoryID, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("error while reading products: %w", err)
		}
		products = append(products, &p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while processing products: %w", err)
	}
	return products, nil
}

func (r *postgresProductRepository) ListByCategory(ctx context.Context, categoryID int64, limit, offset int) ([]*model.Product, error) {
	query := `
		SELECT id, name, description, price, category_id, created_at, updated_at
		FROM products
		WHERE category_id = $1
		ORDER BY id
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(ctx, query, categoryID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error while retrieving products by category: %w", err)
	}
	defer rows.Close()

	products := make([]*model.Product, 0)
	for rows.Next() {
		var p model.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.CategoryID, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("error while reding products: %w", err)
		}
		products = append(products, &p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while processing product lines: %w", err)
	}
	return products, nil
}
