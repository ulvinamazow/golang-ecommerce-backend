package service

import "errors"

var (
	ErrInvalidCredentials = errors.New("email or password is incorrect")
)
