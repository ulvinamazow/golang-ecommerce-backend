package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ulvinamazow/go-backend-project/internal/model"
)

type postgresUserRepository struct {
	db *pgxpool.Pool
}

func NewPostgresUserRepository(db *pgxpool.Pool) UserRepository {
	return &postgresUserRepository{db: db}
}

func (r *postgresUserRepository) Create(ctx context.Context, user *model.User) error {
	query := `
	INSERT INTO users (name, email, password_hash, created_at, updated_at)
	VALUES ($1, $2, NOW(), NOW())
	RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query, user.Email, user.Password).
		Scan(&user.Name, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return fmt.Errorf("User could not be created: %w", err)
	}

	fmt.Println("User created")

	return nil
}

func (r *postgresUserRepository) GetByID(ctx context.Context, id int64) (*model.User, error) {
	query := `
		SELECT id, email, password_hash, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	var user model.User

	err := r.db.QueryRow(ctx, query, id).
		Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("istifadəçi gətirilə bilmədi: %w", err)
	}

	return &user, nil
}

func (r *postgresUserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT id, email, password_hash, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var user model.User

	err := r.db.QueryRow(ctx, query, email).
		Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("user could not be retrieved using the email: %w", err)
	}

	return &user, nil
}

func (r *postgresUserRepository) Update(ctx context.Context, user *model.User) error {
	query := `
		UPDATE users
		SET email = $1, password_hash = $2, updated_at = NOW()
		WHERE id = $3
		RETURNING updated_at
	`

	err := r.db.QueryRow(ctx, query, user.Email, user.Password, user.ID).
		Scan(&user.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return fmt.Errorf("istifadəçi yenilənə bilmədi: %w", err)
	}

	return nil
}

func (r *postgresUserRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("istifadəçi silinə bilmədi: %w", err)
	}

	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *postgresUserRepository) List(ctx context.Context, limit, offset int) ([]*model.User, error) {
	query := `
		SELECT id, email, password_hash, created_at, updated_at
		FROM users
		ORDER BY id
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)

	if err != nil {
		return nil, fmt.Errorf("The list of users could not be retrieved: %w", err)
	}

	defer rows.Close()

	users := make([]*model.User, 0)

	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, fmt.Errorf("Error reading user array: %w", err)
		}
		users = append(users, &u)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error while processing user arrays: %w", err)
	}

	return users, nil
}
