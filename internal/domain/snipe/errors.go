package snipe

import (
	"errors"
)

var (
	ErrSnipeNotFound  = errors.New("snipe not found")
	ErrSnipeDuplicate = errors.New("snipe duplicate")
)
