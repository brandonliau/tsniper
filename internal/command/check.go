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

type checkCommand struct {
	dCfg  *config.DiscordConfig
	repo  repository.Repository
	codec codec.Codec
}

func NewCheckCommand(dCfg *config.DiscordConfig, repo repository.Repository, codec codec.Codec) *checkCommand {
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
		rsp := shared.EphemeralEmbedResponse(shared.InvalidCheck())
		return rsp, nil
	}

	// hash and check if hash already exists in database
	hash := c.codec.Hash(snipes)
	existingChunks, err := c.repo.RetrievePaginationEntry(hash)
	texts := make([]string, 0)

	// if there was an error and it wasn't and no row error
	if err != nil && err != sql.ErrNoRows {
		rsp := shared.EphemeralContentResponse("Something went wrong!")
		return rsp, err
	}

	// check if chunks already exist
	var builder strings.Builder
	if existingChunks != "" {
		err = json.Unmarshal([]byte(existingChunks), &texts)
		if err != nil {
			rsp := shared.EphemeralContentResponse("Something went wrong!")
			return rsp, err
		}
	} else {
		chunked := shared.Chunk(snipes, shared.SnipesPerPage)
		for _, chunk := range chunked {
			builder.Reset()
			for _, snipe := range chunk {
				index, campus, season := snipe[0], snipe[1], snipe[2]
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
			texts = append(texts, builder.String())
		}
		data, err := json.Marshal(texts)
		if err != nil {
			rsp := shared.EphemeralContentResponse("Something went wrong!")
			return rsp, err
		}
		c.repo.AddPaginationEntry(hash, data, time.Now().Unix())
	}

	// create pagination buttons
	var buttons []discordgo.MessageComponent
	var embed *discordgo.MessageEmbed
	if len(texts) > 1 {
		moveButtons := c.repo.RetrieveComponents("backwardSkip", "previousPage", "nextPage", "forwardSkip")
		pageButton := component.NewPageButton(1, len(texts), hash).Component()
		buttons = []discordgo.MessageComponent{moveButtons[0], moveButtons[1], pageButton, moveButtons[2], moveButtons[3]}
	}

	embed = shared.SuccessfulCheck(texts[0])
	rsp := shared.EphemeralEmbedResponse(embed)
	if buttons != nil {
		shared.AddComponent(rsp, buttons...)
	}
	return rsp, nil
}
