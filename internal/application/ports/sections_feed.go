package ports

import (
	"tsniper/internal/domain/scope"
)

type SectionsFeed interface {
	FetchOpenSections(scope scope.AcademicScope) ([]string, error)
}
