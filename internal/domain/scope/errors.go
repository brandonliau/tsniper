package scope

import (
	"errors"
)

var (
	ErrScopeInvalid  = errors.New("invalid scope")
	ErrCampusInvalid = errors.New("invalid campus")
	ErrSeasonInvalid = errors.New("invalid season")
	ErrTermInvalid   = errors.New("invalid term")
)
