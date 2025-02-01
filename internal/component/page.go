package component

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type pageButton struct {
	currentPage int
	totalPages  int
	hash        string
}

func NewPageButton(currentPage int, totalPages int, hash string) *pageButton {
	return &pageButton{
		currentPage: currentPage,
		totalPages:  totalPages,
		hash:        hash,
	}
}

func (c *pageButton) Component() discordgo.MessageComponent {
	return discordgo.Button{
		Label:    fmt.Sprintf("%d/%d", c.currentPage, c.totalPages),
		Style:    discordgo.SecondaryButton,
		CustomID: c.hash,
		Disabled: true,
	}
}
