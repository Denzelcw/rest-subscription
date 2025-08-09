package storage

import "errors"

var (
	ErrNotFound     = errors.New("user_subscription not found")
	ErrUrlExists    = errors.New("user_subscription already exists")
	ErrUserNotFound = errors.New("user not found")
)
