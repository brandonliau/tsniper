package view

import (
	"sort"

	"tsniper/internal/domain/course"
)

type CourseView struct {
	Index        string
	Title        string
	CourseString string
	Section      string
	Instructors  string
	Notes        string
	Meeting      string
	Campus       string
	Term         string
	Year         string
	LastOpen     int64
}

func FromCourse(crs *course.Course) *CourseView {
	return &CourseView{
		Index:        crs.Index,
		Title:        crs.Title,
		CourseString: crs.CourseString,
		Section:      crs.Section,
		Instructors:  crs.Instructors,
		Notes:        crs.Notes,
		Meeting:      crs.Meeting,
		Campus:       string(crs.Scope.Campus),
		Term:         string(crs.Scope.Term),
		Year:         crs.Scope.Year,
		LastOpen:     crs.LastOpen,
	}
}

func FromCourses(courses []*course.Course) []*CourseView {
	out := make([]*CourseView, 0, len(courses))
	for _, c := range courses {
		out = append(out, FromCourse(c))
	}
	return out
}

func SortCourseViews(courses []*CourseView, counts []int) {
	indices := make([]int, len(courses))
	for i := range indices {
		indices[i] = i
	}
	sort.SliceStable(indices, func(i, j int) bool {
		a, b := courses[indices[i]], courses[indices[j]]
		if a.Campus != b.Campus {
			return a.Campus < b.Campus
		}
		if a.Term != b.Term {
			return a.Term < b.Term
		}
		if a.Year != b.Year {
			return a.Year < b.Year
		}
		return a.Index < b.Index
	})
	sortedCourses := make([]*CourseView, len(courses))
	sortedCounts := make([]int, len(counts))
	for i, idx := range indices {
		sortedCourses[i] = courses[idx]
		sortedCounts[i] = counts[idx]
	}
	copy(courses, sortedCourses)
	copy(counts, sortedCounts)
}
