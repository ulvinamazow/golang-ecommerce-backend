package middleware

import "errors"

var (
	ErrUserIDNotFound  = errors.New("context error")
	ErrUserIDWrongType = errors.New("user_id wrong type")
	ErrUserIDEmpty     = errors.New("user_id is empty")
	ErrRoleNotFound    = errors.New("role is not in the context")
	ErrRoleWrongType   = errors.New("role wrong type")
)
