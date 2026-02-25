package command

import (
	"tsniper/internal/shared"

	"github.com/bwmarrin/discordgo"
)

var CampusName = map[string]string{
	"NB": "new brunswick",
	"NK": "newark",
	"CM": "camden",
}

type Command interface {
	Command() *discordgo.ApplicationCommand
	Execute(args *shared.CmdArgs) (*discordgo.InteractionResponse, error)
}

func ParseInteractionOptions(cid discordgo.ApplicationCommandInteractionData) map[string]string {
	opts := make(map[string]string)
	for _, opt := range cid.Options {
		opts[opt.Name] = opt.Value.(string)
	}
	return opts
}
