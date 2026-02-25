package component

import (
	"fmt"
	"net/url"

	"tsniper/internal/repository"
	"tsniper/internal/shared"
	"tsniper/pkg/config"
	"tsniper/pkg/database"

	"github.com/bwmarrin/discordgo"
)

type resnipeButton struct {
	sCfg *config.ServiceConfig
	repo repository.Repository
	db   database.Database
}

func NewResnipeButton(sCfg *config.ServiceConfig, repo repository.Repository, db database.Database) *resnipeButton {
	return &resnipeButton{
		sCfg: sCfg,
		repo: repo,
		db:   db,
	}
}

func (c *resnipeButton) CustomID() string {
	return c.Component().(discordgo.Button).CustomID
}

func (c *resnipeButton) Component() discordgo.MessageComponent {
	return discordgo.Button{
		Label:    "Resnipe",
		Style:    discordgo.PrimaryButton,
		CustomID: "resnipe",
	}
}

func (c *resnipeButton) getSeason(season string) string {
	// ex. input = 12025 (term + season)
	for k, v := range c.sCfg.SeasonData {
		if v.Term+v.Year == season {
			return k
		}
	}
	return ""
}

func (c *resnipeButton) Execute(args *shared.CmdArgs) (*discordgo.InteractionResponse, error) {
	components := args.Interaction.Message.Components[0]
	actionsRow := components.(*discordgo.ActionsRow)
	button := actionsRow.Components[0].(*discordgo.Button)

	// parse index, campus, and season from register button
	resnipeURL, _ := url.Parse(button.URL)
	index := resnipeURL.Query().Get("indexList")
	campus := resnipeURL.Query().Get("campus")
	season := c.getSeason(resnipeURL.Query().Get("semesterSelection"))

	_, ok := c.sCfg.ValidSeasons[season]
	if !ok {
		rsp := shared.EphemeralContentResponse("Cannot re-add snipe from inactive season.")
		return rsp, nil
	}

	if c.repo.IsSniping(args.UserID, index, campus, season) {
		rsp := shared.EphemeralContentResponse(fmt.Sprintf("You are already sniping `%s`.", index))
		return rsp, nil
	}

	c.repo.AddSnipe(args.UserID, index, campus, season)
	c.repo.Add(index, campus, season)

	rsp := shared.EphemeralContentResponse(fmt.Sprintf("Successfully re-added `%s` to your snipe requests.", index))
	return rsp, nil
}
