package authservice

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/ulvinamazow/go-backend-project/internal/auth"
	"github.com/ulvinamazow/go-backend-project/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type authServiceImpl struct {
	userRepo repository.UserRepository
	tokenCfg auth.TokenConfig
}

func NewAuthService(
	userRepo repository.UserRepository,
	tokenCfg auth.TokenConfig,
) *authServiceImpl {
	return &authServiceImpl{
		userRepo: userRepo,
		tokenCfg: tokenCfg,
	}
}

func (authserviceImpl *authServiceImpl) Login(ctx context.Context, email, password string) (string, error) {
	user, err := authserviceImpl.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return "", ErrInvalidCredentials
		}
		return "", fmt.Errorf("error while login: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	); err != nil {
		return "", ErrInvalidCredentials
	}

	userID := strconv.Itoa(int(user.ID))

	token, err := auth.GenerateAccessToken(authserviceImpl.tokenCfg, userID, user.Role)

	if err != nil {
		return "", fmt.Errorf("token could not be created: %w", err)
	}

	return token, nil
}
