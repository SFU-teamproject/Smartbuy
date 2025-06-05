package apperrors

import (
	"errors"
)

var (
	ErrBadRequest         = errors.New("bad request")
	ErrInternal           = errors.New("internal server error")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
	ErrNotFound           = errors.New("not found")
	ErrAlreadyExists      = errors.New("already exists")
)
