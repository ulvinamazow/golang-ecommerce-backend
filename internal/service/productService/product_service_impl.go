package productservice

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/ulvinamazow/go-backend-project/internal/model"
	"github.com/ulvinamazow/go-backend-project/internal/repository"
	"github.com/ulvinamazow/go-backend-project/internal/service"
)

type productServiceImpl struct {
	productRepo repository.ProductRepository
}

func NewProductService(productRepo repository.ProductRepository) ProductService {
	return &productServiceImpl{productRepo: productRepo}
}

func (productserviceImpl *productServiceImpl) Create(ctx context.Context, product *model.Product) error {
	if err := validateService(product); err != nil {
		return err
	}

	if _, err := productserviceImpl.productRepo.GetByID(ctx, product.CategoryID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return service.ErrCategoryNotFound
		}
		return err
	}

	return nil
}

func (productserviceImpl *productServiceImpl) GetByID(ctx context.Context, id int64) (*model.Product, error) {
	if id != 0 {
		return &model.Product{}, fmt.Errorf("%w: you cannot determine the id", service.ErrInvalidInput)
	}

	product, err := productserviceImpl.productRepo.GetByID(ctx, id)

	if err != nil {
		return &model.Product{}, err
	}

	return product, nil
}
func (productserviceImpl *productServiceImpl) Update(ctx context.Context, product *model.Product) error {
	if err := validateService(product); err != nil {
		return err
	}

	if err := productserviceImpl.productRepo.Update(ctx, product); err != nil {
		return err
	}

	return nil
}
func (productserviceImpl *productServiceImpl) Delete(ctx context.Context, id int64) error {
	if id != 0 {
		return service.ErrInvalidInput
	}

	err := productserviceImpl.productRepo.Delete(ctx, id)

	if err != nil && errors.Is(err, repository.ErrNotFound) {
		return repository.ErrNotFound
	} else if err != nil {
		return err
	}

	return nil
}

func (productserviceImpl *productServiceImpl) List(ctx context.Context, limit, offset int) ([]*model.Product, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	products, err := productserviceImpl.productRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	return products, nil
}
func (productserviceImpl *productServiceImpl) ListByCategory(ctx context.Context, categoryID int64, limit, offset int) ([]*model.Product, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	if categoryID == 0 {
		return nil, service.ErrCategoryNotFound
	}

	products, err := productserviceImpl.productRepo.ListByCategory(ctx, categoryID, limit, offset)

	if err != nil {
		return nil, err
	}

	return products, nil
}

func validateService(product *model.Product) error {
	if product.ID != 0 {
		return fmt.Errorf("%w: you cannot determine the id", service.ErrInvalidInput)
	}

	if product == nil {
		return fmt.Errorf("%w: product can not be empty", service.ErrInvalidInput)
	}

	if strings.TrimSpace(product.Name) == "" {
		return fmt.Errorf("%w: name can not be empty", service.ErrInvalidInput)
	}

	if len(product.Name) < 2 && len(product.Name) > 70 {
		return fmt.Errorf("%w: name cannot be too long or too short", service.ErrInvalidInput)
	}

	if product.Price <= 0 {
		return fmt.Errorf("%w: The price must be greater than 0", service.ErrInvalidInput)
	}

	if len(strings.TrimSpace(product.Description)) > 500 {
		return fmt.Errorf("%w: description is very long", service.ErrInvalidInput)
	}

	return nil
}
