package user

import (
	"errors"
)

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrUserDuplicate = errors.New("user duplicate")
)
