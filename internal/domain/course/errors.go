package course

import (
	"errors"
)

var (
	ErrCourseNotFound  = errors.New("course not found")
	ErrCourseDuplicate = errors.New("course duplicate")
)
