package snipe

import (
	"tsniper/internal/domain/scope"
)

type SnipeRepository interface {
	Create(snp *Snipe) error
	Delete(snp *Snipe) error

	Get(userID string, index string, scope scope.AcademicScope) (*Snipe, error)
	GetAll() ([]*Snipe, error)
	GetByUser(userID string) ([]*Snipe, error)
	GetByIndex(index string, scope scope.AcademicScope) ([]*Snipe, error)

	DeleteByUser(userID string) error
}
