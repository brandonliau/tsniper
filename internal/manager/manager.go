package manager

import (
	"tsniper/internal/command"
	"tsniper/internal/component"
)

type Manager interface {
	RegisterCommand(c command.Command)
	RegisterComponent(c component.Component)
}
