package snipe

import (
	"tsniper/internal/domain/scope"
)

type SnipeRepository interface {
	Create(snp *Snipe) error
	Delete(snp *Snipe) error

	Get(userID string, index string, scope scope.AcademicScope) (*Snipe, error)
	ListAll() ([]*Snipe, error)
	ListByUser(userID string) ([]*Snipe, error)
	ListByIndex(index string, scope scope.AcademicScope) ([]*Snipe, error)

	DeleteByUser(userID string) error
}
