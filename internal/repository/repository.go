package repository

import (
	"Tsniper/internal/shared"
)

type Repository interface {
	// registered commands
	Register(name, id string)
	Registered() map[string]string

	// in-memory repository
	TrackedIndices(campus, season string) []string
	Sync()
	Add(index, campus, season string)
	Remove(index, campus, season string)

	// db user management
	Snipes(userID string) [][]string // [][]string{index, campus, season}
	IsSniping(userID, index, campus, season string) bool

	// db snipe management
	AddSnipe(userID, index, campus, season string)
	RemoveSnipe(userID, index, campus, season string)
	ClearSnipe(userID string)

	// db course management
	CourseEntry(index, campus, season string) shared.CourseEntry
	Exists(index, campus, season string) bool
	SnipeUsers(index, campus, season string) []string
	SnipeCount(index, campus, season string) int
	LastOpen(index, campus, season string) int
	LastOpens(campus, season string) map[string]int64
	UpdateLastOpen(index, campus, season string, lastOpen int64)

	// discord user management
	Campus(userID string) string
}
