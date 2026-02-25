package ports

import (
	"tsniper/internal/domain/course"
	"tsniper/internal/domain/scope"
)

type CourseFeed interface {
	FetchCourses(scope scope.AcademicScope) ([]*course.Course, error)
}
