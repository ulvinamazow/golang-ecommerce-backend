package categoryservice

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ulvinamazow/go-backend-project/internal/model"
	categoryrepository "github.com/ulvinamazow/go-backend-project/internal/repository/categoryRepository"
	"github.com/ulvinamazow/go-backend-project/internal/service"
)

type categoryServiceImpl struct {
	categoryRepo categoryrepository.CategoryRepository
}

func NewCategoryService(categoryRepo categoryrepository.CategoryRepository) CategoryService {
	return &categoryServiceImpl{categoryRepo: categoryRepo}
}

func (categoryserviceImpl *categoryServiceImpl) Create(ctx context.Context, category *model.Category) error {
	err := validateCategory(category)

	if err != nil {
		return fmt.Errorf("%w: %w", service.ErrInvalidInput, err)
	}

	category.CreatedAt = time.Now()

	category.UpdatedAt = time.Now()

	categoryserviceImpl.categoryRepo.Create(ctx, category)

	return nil
}

func (categoryserviceImpl *categoryServiceImpl) GetByID(ctx context.Context, id int64) (*model.Category, error) {
	if id <= 0 {
		return &model.Category{}, fmt.Errorf("%w: You cannot determine the id", service.ErrInvalidInput)
	}

	category, err := categoryserviceImpl.categoryRepo.GetByID(ctx, id)

	if err != nil {
		return &model.Category{}, err
	}

	return category, nil
}

func (categoryserviceImpl *categoryServiceImpl) Update(ctx context.Context, category *model.Category) error {
	err := validateCategory(category)
	if err != nil {
		return err
	}

	if err = categoryserviceImpl.categoryRepo.Update(ctx, category); err != nil {
		return err
	}

	return nil
}

func (categoryserviceImpl *categoryServiceImpl) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("%w: You cannot determine the id", service.ErrInvalidInput)
	}

	err := categoryserviceImpl.categoryRepo.Delete(ctx, id)

	if err != nil {
		return err
	}

	return nil
}

func (categoryserviceImpl *categoryServiceImpl) List(ctx context.Context, limit, offset int) ([]*model.Category, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	categories, err := categoryserviceImpl.categoryRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func validateCategory(category *model.Category) error {
	if strings.TrimSpace(category.Name) == "" {
		return fmt.Errorf("%w: category name can not be empty", service.ErrInvalidInput)
	}

	if len(category.Name) < 2 && len(category.Name) > 30 {
		return fmt.Errorf("%w: category name is too long or too short", service.ErrInvalidInput)
	}

	if category.ID != 0 {
		return fmt.Errorf("%w: you cannot determine the id", service.ErrInvalidInput)
	}

	return nil
}
