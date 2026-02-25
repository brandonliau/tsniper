package course

import (
	"tsniper/internal/domain/scope"
)

type Course struct {
	Index        string
	Title        string
	CourseString string
	Section      string
	Instructors  string
	Notes        string
	Meeting      string
	Scope        scope.AcademicScope
}

func NewCourse(index string, title string, courseString string, section string, instructors string, notes string, meeting string, scope scope.AcademicScope) *Course {
	return &Course{
		Index:        index,
		Title:        title,
		CourseString: courseString,
		Section:      section,
		Instructors:  instructors,
		Notes:        notes,
		Meeting:      meeting,
		Scope:        scope,
	}
}
