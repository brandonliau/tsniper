package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func SendEmbedReponse(s *discordgo.Session, i *discordgo.InteractionCreate, embed *discordgo.MessageEmbed) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse {
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData {
			Flags: discordgo.MessageFlagsEphemeral,
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}

func SendContentMessage(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse {
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData {
			Flags: discordgo.MessageFlagsEphemeral,
			Content: content,
		},
	})
}

func SendComplexMessage(user *discordgo.User, dmChannel *discordgo.Channel, embed *discordgo.MessageEmbed, buttons ...discordgo.Button) error {
	components := make([]discordgo.MessageComponent, len(buttons))
	for i, button := range buttons {
		components[i] = button
	}
	_, err := s.ChannelMessageSendComplex(
		dmChannel.ID,
		&discordgo.MessageSend{
			Content: user.Mention(),
			Embeds: []*discordgo.MessageEmbed{embed},
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: components,
				},
			},
		},
	)
	return err
}

func (runTimeData RunTimeData) CreateRegButton(course CourseData, season Season, campus string) discordgo.Button {
	return discordgo.Button{
		Label: "Register",
		Style: discordgo.LinkButton,
		URL: fmt.Sprintf(
			"https://sims.rutgers.edu/webreg/editSchedule.htm?login=cas&semesterSelection=%s&indexList=%s&campus=%s",
			season.Term + season.Year,
			course.Index,
			campus,
		),
	}
}

func CreateResnipeButton() discordgo.Button {
	return discordgo.Button{
		Label: "Resnipe",
		Style: discordgo.PrimaryButton,
		CustomID: "resnipe",
	}
}
