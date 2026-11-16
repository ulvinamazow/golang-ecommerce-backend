package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	apperrors "github.com/ulvinamazow/go-backend-project/internal/errors"
	"github.com/ulvinamazow/go-backend-project/internal/service"
	authservice "github.com/ulvinamazow/go-backend-project/internal/service/authService"
)

var (
	ErrInvalidInput = errors.New("invalid input")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrNotFound     = errors.New("not found")
	ErrInternal     = errors.New("internal server error")
)

const (
	MsgCreated = "product created"
	MsgUpdated = "product updated"
	MsgDeleted = "product deleted"
)

func RespondError(c *gin.Context, err error) error {
	switch {
	case errors.Is(err, ErrInvalidInput), errors.Is(err, service.ErrInvalidInput), errors.Is(err, apperrors.ErrInvalidReference):
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidInput.Error()})
	case errors.Is(err, ErrUnauthorized), errors.Is(err, authservice.ErrInvalidCredentials):
		c.JSON(http.StatusUnauthorized, gin.H{"error": ErrUnauthorized.Error()})
	case errors.Is(err, ErrForbidden), errors.Is(err, service.ErrForbidden):
		c.JSON(http.StatusForbidden, gin.H{"error": ErrForbidden.Error()})
	case errors.Is(err, ErrNotFound), errors.Is(err, apperrors.ErrNotFound), errors.Is(err, service.ErrCategoryNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": ErrNotFound.Error()})
	case errors.Is(err, apperrors.ErrAlreadyExists), errors.Is(err, service.ErrEmailAlreadyExists):
		c.JSON(http.StatusConflict, gin.H{"error": apperrors.ErrAlreadyExists.Error()})
	default:
		slog.Error("request failed", "path", c.FullPath(), "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrInternal.Error()})
	}
	return err
}
