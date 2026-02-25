package command

import (
	"fmt"
	"time"

	"tsniper/internal/config"
	"tsniper/internal/interfaces/discord/interaction"
	"tsniper/internal/interfaces/discord/presentation"

	"github.com/bwmarrin/discordgo"
)

type helpCommand struct {
	applicationID string
	customization *config.CustomizationConfig
}

func HelpCommandDefinition() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "help",
		Description: "List all commands.",
	}
}

func HelpCommandHandler(applicationID string, customization *config.CustomizationConfig) interaction.HandleFunc {
	c := &helpCommand{
		applicationID: applicationID,
		customization: customization,
	}
	return c.execute
}

func (c *helpCommand) execute(s *discordgo.Session, i *discordgo.InteractionCreate) (*discordgo.InteractionResponse, error) {
	commands, err := s.ApplicationCommands(c.applicationID, "")
	if err != nil {
		return nil, err
	}

	rsp := interaction.InteractionInitialResponse(
		interaction.WithEmbeds(c.helpEmbed(commands)),
		interaction.WithEphemeral(),
	)
	return rsp, nil
}

func (c *helpCommand) helpEmbed(commands []*discordgo.ApplicationCommand) *discordgo.MessageEmbed {
	commandDefinitions := make(map[string]*discordgo.ApplicationCommand)
	for _, cmd := range commands {
		commandDefinitions[cmd.Name] = cmd
	}
	format := func(def *discordgo.ApplicationCommand) string {
		cmd, ok := commandDefinitions[def.Name]
		if !ok || cmd.ID == "" {
			return fmt.Sprintf("`/%s` - %s", def.Name, def.Description)
		}
		return fmt.Sprintf("</%s:%s> - %s", cmd.Name, cmd.ID, cmd.Description)
	}
	return &discordgo.MessageEmbed{
		Title: "Commands",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name: "Sniping",
				Value: fmt.Sprintf(
					"%s\n%s\n%s\n%s\n%s",
					format(AddCommandDefinition()),
					format(RemoveCommandDefinition()),
					format(ClearCommandDefinition()),
					format(CheckCommandDefinition()),
					format(SearchCommandDefinition()),
				),
				Inline: false,
			},
			{
				Name: "Miscellaneous",
				Value: fmt.Sprintf(
					"%s",
					format(StatusCommandDefinition()),
				),
				Inline: false,
			},
		},
		Color: presentation.Blue,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: c.customization.Thumbnail,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: time.Now().Format("01/02/2006 03:04:05 PM"),
		},
	}
}
