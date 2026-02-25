package snipe

import (
	"tsniper/internal/domain/scope"
)

type Snipe struct {
	UserID string
	Index  string
	Scope  scope.AcademicScope
}

func NewSnipe(userID string, index string, scope scope.AcademicScope) *Snipe {
	return &Snipe{
		UserID: userID,
		Index:  index,
		Scope:  scope,
	}
}
