package course

import (
	"tsniper/internal/domain/scope"
)

type CourseRepository interface {
	Create(crs *Course) error
	Save(crs *Course) error
	Delete(crs *Course) error

	Get(index string, scope scope.AcademicScope) (*Course, error)
	GetAll() ([]*Course, error)
	GetAllByScope(scope scope.AcademicScope) ([]*Course, error)

	BatchCreate(courses []*Course) error
}
