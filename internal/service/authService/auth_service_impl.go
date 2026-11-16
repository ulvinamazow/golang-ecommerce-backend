package authservice

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/ulvinamazow/go-backend-project/internal/auth"
	// OLD: "github.com/ulvinamazow/go-backend-project/internal/repository"
	"github.com/ulvinamazow/go-backend-project/internal/errors"
	userrepository "github.com/ulvinamazow/go-backend-project/internal/repository/userRepository"
	"golang.org/x/crypto/bcrypt"
)

type authServiceImpl struct {
	userRepo userrepository.UserRepository
	tokenCfg auth.TokenConfig
}

func NewAuthService(
	userRepo userrepository.UserRepository,
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
		if errors.Is(err, apperrors.ErrNotFound) {
			// OLD: return "", ErrInvalidCredentials
			return "", fmt.Errorf("login user: %w: %w", ErrInvalidCredentials, err)
		}
		// OLD: return "", fmt.Errorf("error while login: %w", err)
		return "", fmt.Errorf("login user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	); err != nil {
		// OLD: return "", ErrInvalidCredentials
		return "", fmt.Errorf("login user: %w: %w", ErrInvalidCredentials, err)
	}

	userID := strconv.Itoa(int(user.ID))

	token, err := auth.GenerateAccessToken(authserviceImpl.tokenCfg, userID, user.Role)

	if err != nil {
		// OLD: return "", fmt.Errorf("token could not be created: %w", err)
		return "", fmt.Errorf("create access token: %w", err)
	}

	return token, nil
}
