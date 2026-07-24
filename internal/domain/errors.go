package domain

import "errors"

var (
	ErrNotFound          = errors.New("not found")
	ErrConflict          = errors.New("already exists")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrForbidden         = errors.New("forbidden")
	ErrValidation        = errors.New("validation error")
	ErrInternal          = errors.New("internal error")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrTokenExpired      = errors.New("token expired")
	ErrTokenInvalid      = errors.New("token invalid")
)
