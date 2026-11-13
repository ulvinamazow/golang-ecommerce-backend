package productrepository

import (
	"context"

	"github.com/ulvinamazow/go-backend-project/internal/model"
)

// ProductRepository məhsul məlumatları üzərində əməliyyatları təyin edir.
type ProductRepository interface {
	Create(ctx context.Context, product *model.Product) error
	GetByID(ctx context.Context, id int64) (*model.Product, error)
	Update(ctx context.Context, product *model.Product) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, limit, offset int) ([]*model.Product, error)
	ListByCategory(ctx context.Context, categoryID int64, limit, offset int) ([]*model.Product, error)
}
