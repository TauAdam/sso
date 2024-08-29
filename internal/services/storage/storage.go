package storage

import "errors"

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrUserDuplicate = errors.New("user already exists")
	ErrAppNotFound   = errors.New("app not found")
	ErrAppDuplicate  = errors.New("app already exists")
)
