package storage

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrEmailExists  = errors.New("email exists")
)
