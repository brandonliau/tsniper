package repository

import (
	"tsniper/internal/shared"

	"github.com/bwmarrin/discordgo"
)

type Repository interface {
	// command registration
	RegisterCommand(name string, id string)
	RetrieveCommands(names ...string) []string

	// component registration
	RegisterComponent(name string, component discordgo.MessageComponent)
	RetrieveComponents(names ...string) []discordgo.MessageComponent

	// in-memory snipe management
	TrackedIndices(campus, season string) []string
	Sync()
	Add(index, campus, season string)
	Remove(index, campus, season string)

	// db user management
	Snipes(userID string) []shared.Snipe
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

	// db pagination management
	AddPaginationEntry(hash string, data []byte, datetime int64)
	RetrievePaginationEntry(hash string) ([]byte, error)

	// discord user management
	Campus(userID string) string
}
