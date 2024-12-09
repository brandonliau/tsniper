package component

import (
	"fmt"
	"net/url"
	"time"

	"Tsniper/internal/repository"
	"Tsniper/internal/shared"
	"Tsniper/pkg/config"
	"Tsniper/pkg/database"

	"github.com/bwmarrin/discordgo"
)

type resnipeButton struct {
	dCfg *config.DiscordConfig
	sCfg *config.ServiceConfig
	repo repository.Repository
	db   database.Database
}

func NewResnipeButton(dCfg *config.DiscordConfig, sCfg *config.ServiceConfig, repo repository.Repository, db database.Database) *resnipeButton {
	return &resnipeButton{
		dCfg: dCfg,
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

	if c.repo.IsSniping(args.UserID, index, campus, season) {
		rsp := shared.EphemeralContentResponse(fmt.Sprintf("You are already sniping `%s`.", index))
		return rsp, nil
	}

	c.repo.AddSnipe(args.UserID, index, campus, season)
	c.repo.Add(index, campus, season)

	course := c.repo.CourseEntry(index, campus, season)
	rsp := shared.EphemeralEmbedResponse(c.SuccessfullReAdd(course))
	return rsp, nil
}

func (c *resnipeButton) SuccessfullReAdd(course shared.CourseEntry) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title: "Successfully Added Request!",
		Description: fmt.Sprintf(
			"`%s` - %s (**Section %s**) was added to your snipe requests.",
			course.Index,
			course.Title,
			course.Section,
		),
		Color: shared.Green,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: c.dCfg.Image,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: time.Now().Format("01/02/2006 03:04:05 PM"),
		},
	}
}
