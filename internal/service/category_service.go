package service

import (
	"context"

	"github.com/ulvinamazow/go-backend-project/internal/model"
)

// CategoryService kateqoriya ilə bağlı biznes əməliyyatlarını müəyyən edir.
type CategoryService interface {
	Create(ctx context.Context, category *model.Category) error
	GetByID(ctx context.Context, id int64) (*model.Category, error)
	Update(ctx context.Context, category *model.Category) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, limit, offset int) ([]*model.Category, error)
}
