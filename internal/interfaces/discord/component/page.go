package component

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func PageComponentDefinition(currentPage int, totalPages int) discordgo.Button {
	return discordgo.Button{
		CustomID: "page",
		Label:    fmt.Sprintf("Page %d / %d", currentPage, totalPages),
		Style:    discordgo.SecondaryButton,
		Disabled: true,
	}
}
