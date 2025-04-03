package command

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"Tsniper/internal/component"
	"Tsniper/internal/repository"
	"Tsniper/internal/shared"

	"Tsniper/pkg/codec"
	"Tsniper/pkg/config"

	"github.com/bwmarrin/discordgo"
)

const (
	snipesPerPage = 30
)

type checkCommand struct {
	dCfg  *config.DiscordConfig
	repo  repository.Repository
	codec codec.Codec[shared.Snipe]
}

func NewCheckCommand(dCfg *config.DiscordConfig, repo repository.Repository, codec codec.Codec[shared.Snipe]) *checkCommand {
	return &checkCommand{
		dCfg:  dCfg,
		repo:  repo,
		codec: codec,
	}
}

func (c *checkCommand) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "check",
		Description: "View all active snipe requests.",
	}
}

func (c *checkCommand) Execute(args *shared.CmdArgs) (*discordgo.InteractionResponse, error) {
	snipes := c.repo.Snipes(args.UserID)
	if len(snipes) == 0 {
		return shared.EphemeralEmbedResponse(shared.InvalidCheck()), nil
	}

	// hash snipes
	hash, err := c.codec.Hash(snipes)
	if err != nil {
		return shared.EphemeralContentResponse("Something went wrong!"), err
	}

	// attempt to retrieve chunks
	existingChunks, err := c.repo.RetrievePaginationEntry(hash)
	if err != nil && err != sql.ErrNoRows {
		return shared.EphemeralContentResponse("Something went wrong!"), err
	}
	
	// check if chunks already exist
	var stringChunks []string
	var builder strings.Builder
	var buttons []discordgo.MessageComponent
	
	if existingChunks != nil {
		err := json.Unmarshal(existingChunks, &stringChunks)
		if err != nil {
			return shared.EphemeralContentResponse("Something went wrong!"), err
		}
	} else {
		chunked := shared.Chunk(snipes, snipesPerPage)
		for _, chunk := range chunked {
			builder.Reset()
			for _, snipe := range chunk {
				index, campus, season := snipe.Index, snipe.Campus, snipe.Season
				course := c.repo.CourseEntry(index, campus, season)
				count := c.repo.SnipeCount(index, campus, season)
				lastOpen := c.repo.LastOpen(index, campus, season)
				text := fmt.Sprintf("%s `%s` - %s (**Section %s**) | :eyes: `%d` | ", shared.EmojiMap[season], course.Index, course.Title, course.Section, count)
				builder.WriteString(text)
				if lastOpen == -1 {
					builder.WriteString("`Unknown`\n")
				} else {
					builder.WriteString(fmt.Sprintf("<t:%d:R>\n", lastOpen))
				}
			}
			stringChunks = append(stringChunks, builder.String())
		}
	}

	// add pagination entry
	if existingChunks == nil && len(stringChunks) > 1 {
		data, err := json.Marshal(stringChunks)
		if err != nil {
			return shared.EphemeralContentResponse("Something went wrong!"), err
		}
		c.repo.AddPaginationEntry(hash, data, time.Now().Unix())
	}

	// retrieve pagination buttons
	if len(stringChunks) > 1 {
		moveButtons := c.repo.RetrieveComponents("backwardSkip", "previousPage", "nextPage", "forwardSkip")
		pageButton := component.NewPageButton(1, len(stringChunks), hash).Component()
		buttons = []discordgo.MessageComponent{moveButtons[0], moveButtons[1], pageButton, moveButtons[2], moveButtons[3]}
	}

	// create and send embed
	embed := shared.SuccessfulCheck(stringChunks[0])
	rsp := shared.EphemeralEmbedResponse(embed)
	if buttons != nil {
		shared.AddComponent(rsp, buttons...)
	}
	return rsp, nil
}
