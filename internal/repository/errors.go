package repository

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	apperrors "github.com/ulvinamazow/go-backend-project/internal/errors"
)

func MapPgError(op string, err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return fmt.Errorf("%s: %w: %w", op, apperrors.ErrAlreadyExists, err)
		case "23503":
			return fmt.Errorf("%s: %w: %w", op, apperrors.ErrInvalidReference, err)
		}
	}
	return fmt.Errorf("%s: %w", op, err)
}
