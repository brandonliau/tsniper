package snipe

import "tsniper/internal/domain/scope"

type SnipeCache interface {
	Add(snp *Snipe)
	Remove(snp *Snipe)
	Clear(snipes []*Snipe)
	Tracked(scope scope.AcademicScope) []string
}
