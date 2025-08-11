package storage

import "errors"

var (
	ErrNotFound      = errors.New("user_subscription not found")
	ErrUserSubExists = errors.New("user_subscription already exists")
	ErrUserNotFound  = errors.New("user not found")
	ErrOverlap       = errors.New("user subscription conflicts with existing record")
)
