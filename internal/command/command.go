package command

import (
	"Tsniper/internal/shared"

	"github.com/bwmarrin/discordgo"
)

var CampusName = map[string]string{
	"NB": "new brunswick",
	"NK": "newark",
	"CM": "camden",
}

var EmojiMap = map[string]string{
	"spring": ":herb:",
	"summer": ":sunny:",
	"fall":   ":fallen_leaf:",
	"winter": ":snowflake:",
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
