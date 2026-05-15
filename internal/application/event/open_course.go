package event

import (
	"tsniper/internal/application/view"
)

type CourseOpenType int

const (
	CourseOpenNotification CourseOpenType = iota
	CourseOpenAutosnipeSuccess
	CourseOpenAutosnipeFailed
)

// todo: replace []string with []uuid.UUID once there is a way to map external user IDs to internal uuids
type CourseOpen struct {
	Type       CourseOpenType
	ErrMessage string
	Course     *view.CourseView
	UserIDs    []string
}
