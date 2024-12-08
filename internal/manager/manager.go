package manager

import (
	"Tsniper/internal/command"
	"Tsniper/internal/component"
)

type Manager interface {
	RegisterCommand(c command.Command)
	RegisterComponent(c component.Component)
}
