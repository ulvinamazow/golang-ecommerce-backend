package handler

import "errors"

var (
	ErrInvalidMethod = errors.New("The method is incorrect")

	ErrUnauthorized = errors.New("Status unauthorized")

	ErrForbidden = errors.New("Status forbidden")

	ErrInvalidInput = errors.New("Invalid input")

	StatusCreated = errors.New("Status created")

	StatusDeleted = errors.New("Status deleted")

	
	
)