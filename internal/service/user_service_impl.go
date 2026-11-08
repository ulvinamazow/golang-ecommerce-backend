package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/ulvinamazow/go-backend-project/internal/model"
	"github.com/ulvinamazow/go-backend-project/internal/repository"
)

type userServiceImpl struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userServiceImpl{userRepo: userRepo}
}

var _ UserService = (*userServiceImpl)(nil)

func (s *userServiceImpl) Create(ctx context.Context, user *model.User) error {
	if err := validateUser(user); err != nil {
		return err
	}

	user.Email = strings.ToLower(strings.TrimSpace(user.Email))
	existing, err := s.userRepo.GetByEmail(ctx, user.Email)
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return fmt.Errorf("Email could not be verified: %w", err)
	}
	if existing != nil {
		return ErrEmailAlreadyExists
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return fmt.Errorf("istifadəçi yaradıla bilmədi: %w", err)
	}
	return nil
}

func (s *userServiceImpl) GetByID(ctx context.Context, id int64) (*model.User, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: id müsbət olmalidir", ErrInvalidInput)
	}

	return s.userRepo.GetByID(ctx, id)
}

func (s *userServiceImpl) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" {
		return nil, fmt.Errorf("%w: email boş ola bilməz", ErrInvalidInput)
	}
	return s.userRepo.GetByEmail(ctx, email)
}

func (s *userServiceImpl) Update(ctx context.Context, user *model.User) error {
	if err := validateUser(user); err != nil {
		return err
	}
	// Mövcud istifadəçinin olduğunu yoxlayaq
	if _, err := s.userRepo.GetByID(ctx, user.ID); err != nil {
		return err
	}
	user.Email = strings.ToLower(strings.TrimSpace(user.Email))
	return s.userRepo.Update(ctx, user)
}

func (s *userServiceImpl) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("%w: the id ,ust be positive", ErrInvalidInput)
	}
	return s.userRepo.Delete(ctx, id)
}

func (s *userServiceImpl) List(ctx context.Context, limit, offset int) ([]*model.User, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.userRepo.List(ctx, limit, offset)
}

// validateUser biznes səviyyəsində giriş yoxlaması aparır.
func validateUser(user *model.User) error {
	if user == nil {
		return fmt.Errorf("%w: user can not be empty", ErrInvalidInput)
	}
	if strings.TrimSpace(user.Email) == "" {
		return fmt.Errorf("%w: email can not be empty", ErrInvalidInput)
	}
	if !strings.Contains(user.Email, "@") {
		return fmt.Errorf("%w: invalid email format", ErrInvalidInput)
	}
	if user.PasswordHash == "" {
		return fmt.Errorf("%w: password can not be empty", ErrInvalidInput)
	}
	return nil
}
