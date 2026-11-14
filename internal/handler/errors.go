package handler

import "errors"

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
