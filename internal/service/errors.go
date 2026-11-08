package service

import "errors"

var (
	ErrEmailAlreadyExists = errors.New("this email already exists")

	ErrInvalidInput = errors.New("incorrect input data")

	ErrCategoryNotFound = errors.New("The specified category does not exist")
)
