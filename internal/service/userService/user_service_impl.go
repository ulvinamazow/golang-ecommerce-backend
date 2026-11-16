package userservice

import (
	"context"
	"errors"
	"fmt"
	"strings"

	apperrors "github.com/ulvinamazow/go-backend-project/internal/errors"
	"github.com/ulvinamazow/go-backend-project/internal/model"
	userrepository "github.com/ulvinamazow/go-backend-project/internal/repository/userRepository"
	"github.com/ulvinamazow/go-backend-project/internal/service"
	"golang.org/x/crypto/bcrypt"
)

type userServiceImpl struct {
	userRepo userrepository.UserRepository
}

func NewUserService(userRepo userrepository.UserRepository) UserService {
	return &userServiceImpl{userRepo: userRepo}
}

var _ UserService = (*userServiceImpl)(nil)

func (s *userServiceImpl) Create(ctx context.Context, user *model.User) error {
	if err := validateUser(user); err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	hashedPassword, err := generateHash(user.Password)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	user.Password = hashedPassword

	user.Email = strings.TrimSpace(user.Email)

	existing, err := s.userRepo.GetByEmail(ctx, user.Email)

	if err != nil && !errors.Is(err, apperrors.ErrNotFound) {
		return fmt.Errorf("check user email: %w", err)
	}

	if existing != nil {
		return fmt.Errorf("create user: %w", service.ErrEmailAlreadyExists)
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

func (s *userServiceImpl) GetByID(ctx context.Context, id int64) (*model.User, error) {
	if id <= 0 {
		return nil, fmt.Errorf("get user: %w", service.ErrInvalidInput)
	}

	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	return user, nil
}

func (s *userServiceImpl) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" {
		return nil, fmt.Errorf("get user by email: %w", service.ErrInvalidInput)
	}
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("get user by email: %w", err)
	}
	return user, nil
}

func (s *userServiceImpl) Update(ctx context.Context, user *model.User) error {
	if err := validateUser(user); err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	if _, err := s.userRepo.GetByID(ctx, user.ID); err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	user.Email = strings.ToLower(strings.TrimSpace(user.Email))
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	return nil
}

func (s *userServiceImpl) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("delete user: %w", service.ErrInvalidInput)
	}
	if err := s.userRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	return nil
}

func (s *userServiceImpl) List(ctx context.Context, limit, offset int) ([]*model.User, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	users, err := s.userRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	return users, nil
}

func validateUser(user *model.User) error {
	if user == nil {
		return fmt.Errorf("%w: user can not be empty", service.ErrInvalidInput)
	}
	if strings.TrimSpace(user.Email) == "" {
		return fmt.Errorf("%w: email can not be empty", service.ErrInvalidInput)
	}
	if !strings.Contains(user.Email, "@") {
		return fmt.Errorf("%w: invalid email format", service.ErrInvalidInput)
	}
	if user.Password == "" {
		return fmt.Errorf("%w: password can not be empty", service.ErrInvalidInput)
	}
	return nil
}

func generateHash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", fmt.Errorf("create user: %w", err)
	}

	return string(hashedPassword), nil
}
