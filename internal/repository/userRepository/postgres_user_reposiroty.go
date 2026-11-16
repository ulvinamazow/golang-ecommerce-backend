package userrepository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	myerrors "github.com/ulvinamazow/go-backend-project/internal/errors"
	"github.com/ulvinamazow/go-backend-project/internal/model"
	"github.com/ulvinamazow/go-backend-project/internal/repository"
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
	VALUES ($1, $2, $3, NOW(), NOW())
	RETURNING id, name, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query, user.Name, user.Email, user.Password).
		Scan(&user.ID, &user.Name, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		// OLD: return myerrors.DBErrUserCreated
		return repository.MapPgError("create user", err)
	}

	fmt.Println(myerrors.DBUserCreated)

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
			// OLD: return nil, myerrors.DBErrNotFound
			return nil, fmt.Errorf("user %d: %w: %w", id, myerrors.ErrNotFound, err)
		}
		// OLD: return nil, myerrors.DBErrUserNotRetrieved
		return nil, fmt.Errorf("get user: %w", err)
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
			// OLD: return nil, myerrors.DBErrNotFound
			return nil, fmt.Errorf("user with email %s: %w: %w", email, myerrors.ErrNotFound, err)
		}
		// OLD: return nil, myerrors.DBErrUserNotRetrieved
		return nil, fmt.Errorf("get user by email: %w", err)
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
			// OLD: return myerrors.DBErrNotFound
			return fmt.Errorf("user %d: %w: %w", user.ID, myerrors.ErrNotFound, err)
		}
		// OLD: return myerrors.DBErrUserUpdated
		return repository.MapPgError("update user", err)
	}

	return nil
}

func (r *postgresUserRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		// OLD: return fmt.Errorf("user can not deleted: %w", err)
		return repository.MapPgError("delete user", err)
	}

	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		// OLD: return myerrors.DBErrNotFound
		return fmt.Errorf("user %d: %w", id, myerrors.ErrNotFound)
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
		// OLD: return nil, myerrors.DBErrUsersNotRetrieved
		return nil, fmt.Errorf("list users: %w", err)
	}

	defer rows.Close()

	users := make([]*model.User, 0)

	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt); err != nil {
			// OLD: return nil, myerrors.DBErrUsersNotRetrieved
			return nil, fmt.Errorf("list users: %w", err)
		}
		users = append(users, &u)
	}

	if err := rows.Err(); err != nil {
		// OLD: return nil, myerrors.DBErrUsersNotRetrieved
		return nil, fmt.Errorf("list users: %w", err)
	}

	return users, nil
}
