package component

import (
	"errors"
	"fmt"

	"tsniper/internal/application/usecase"
	"tsniper/internal/interfaces/discord/interaction"

	"tsniper/pkg/utils"

	"github.com/bwmarrin/discordgo"
)

type resnipeComponent struct {
	snipeService *usecase.SnipeService
}

func ResnipeComponentDefinition(data ...utils.KeyValue[string, string]) discordgo.Button {
	return discordgo.Button{
		CustomID: interaction.EncodeCustomID("resnipe", data...),
		Label:    "Resnipe",
		Style:    discordgo.PrimaryButton,
	}
}

func ResnipeComponentHandler(snipeService *usecase.SnipeService) interaction.HandleFunc {
	c := &resnipeComponent{
		snipeService: snipeService,
	}
	return c.execute
}

func (c *resnipeComponent) execute(s *discordgo.Session, i *discordgo.InteractionCreate) (*discordgo.InteractionResponse, error) {
	_, params := interaction.DecodeCustomID(i.MessageComponentData().CustomID)

	req := usecase.ReAddSnipeRequest{
		UserID: interaction.GetUserID(i),
		Index:  params["index"],
		Campus: params["campus"],
		Term:   params["term"],
		Year:   params["year"],
	}

	var rsp *discordgo.InteractionResponse
	_, err := c.snipeService.ReAdd(req)
	switch {
	case err == nil:
		rsp = interaction.InteractionInitialResponse(
			interaction.WithContent(fmt.Sprintf("Successfully re-added `%s` to your snipe requests.", req.Index)),
			interaction.WithEphemeral(),
		)
	case errors.Is(err, usecase.ErrReAddSnipeDuplicate):
		rsp = interaction.InteractionInitialResponse(
			interaction.WithContent(fmt.Sprintf("You are already sniping `%s`.", req.Index)),
			interaction.WithEphemeral(),
		)
	case errors.Is(err, usecase.ErrReAddSnipeInvalid):
		rsp = interaction.InteractionInitialResponse(
			interaction.WithContent("Cannot re-add snipe from inactive season."),
			interaction.WithEphemeral(),
		)
	default:
		return nil, err
	}

	return rsp, nil
}
