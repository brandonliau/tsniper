package interaction

import (
	"github.com/bwmarrin/discordgo"
)

func InteractionResponse(responseType discordgo.InteractionResponseType, data *discordgo.InteractionResponseData, opts ...InteractionResponseOption) *discordgo.InteractionResponse {
	for _, opt := range opts {
		opt(data)
	}
	return &discordgo.InteractionResponse{
		Type: responseType,
		Data: data,
	}
}

func InteractionInitialResponse(opts ...InteractionResponseOption) *discordgo.InteractionResponse {
	return InteractionResponse(discordgo.InteractionResponseChannelMessageWithSource, &discordgo.InteractionResponseData{}, opts...)
}

func InteractionUpdateResponse(opts ...InteractionResponseOption) *discordgo.InteractionResponse {
	return InteractionResponse(discordgo.InteractionResponseUpdateMessage, &discordgo.InteractionResponseData{}, opts...)
}

func InteractionModalResponse(data *discordgo.InteractionResponseData) *discordgo.InteractionResponse {
	return InteractionResponse(discordgo.InteractionResponseModal, data)
}

type InteractionResponseOption func(*discordgo.InteractionResponseData)

func WithContent(content string) InteractionResponseOption {
	return func(d *discordgo.InteractionResponseData) {
		d.Content = content
	}
}

func WithEmbeds(embeds ...*discordgo.MessageEmbed) InteractionResponseOption {
	return func(d *discordgo.InteractionResponseData) {
		d.Embeds = embeds
	}
}

func WithComponents(components ...discordgo.MessageComponent) InteractionResponseOption {
	return func(d *discordgo.InteractionResponseData) {
		var rows []discordgo.MessageComponent
		var buttons []discordgo.MessageComponent

		for _, c := range components {
			if _, ok := c.(discordgo.SelectMenu); ok {
				if len(buttons) > 0 {
					rows = append(rows, discordgo.ActionsRow{Components: buttons})
					buttons = nil
				}
				rows = append(rows, discordgo.ActionsRow{Components: []discordgo.MessageComponent{c}})
			} else {
				buttons = append(buttons, c)
				if len(buttons) == 5 {
					rows = append(rows, discordgo.ActionsRow{Components: buttons})
					buttons = nil
				}
			}
		}

		if len(buttons) > 0 {
			rows = append(rows, discordgo.ActionsRow{Components: buttons})
		}

		d.Components = rows
	}
}

func WithFlags(flags discordgo.MessageFlags) InteractionResponseOption {
	return func(d *discordgo.InteractionResponseData) {
		d.Flags = flags
	}
}

func WithEphemeral() InteractionResponseOption {
	return func(d *discordgo.InteractionResponseData) {
		d.Flags = discordgo.MessageFlagsEphemeral
	}
}
