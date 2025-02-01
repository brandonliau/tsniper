package component

import (
	"encoding/json"
	"strconv"

	"Tsniper/internal/repository"
	"Tsniper/internal/shared"

	"Tsniper/pkg/database"

	"github.com/bwmarrin/discordgo"
)

type backwardPageButton struct {
	disabled bool
	repo     repository.Repository
	db       database.Database
}

func NewPreviousPageButton(disabled bool, repo repository.Repository, db database.Database) *backwardPageButton {
	return &backwardPageButton{
		disabled: disabled,
		repo:     repo,
		db:       db,
	}
}

func (c *backwardPageButton) CustomID() string {
	return c.Component().(discordgo.Button).CustomID
}

func (c *backwardPageButton) Component() discordgo.MessageComponent {
	button := discordgo.Button{
		Emoji: &discordgo.ComponentEmoji{
			Name: "⬅",
		},
		Style:    discordgo.PrimaryButton,
		CustomID: "previousPage",
	}
	button.Disabled = c.disabled
	return button
}

func (c *backwardPageButton) Execute(args *shared.CmdArgs) (*discordgo.InteractionResponse, error) {
	embed := args.Interaction.Message.Embeds[0]
	components := args.Interaction.Message.Components[0]
	actionsRow := components.(*discordgo.ActionsRow)

	// get hash from page button and retrieve data
	hash := actionsRow.Components[2].(*discordgo.Button).CustomID
	existingChunks, err := c.repo.RetrievePaginationEntry(hash)
	if err != nil {
		return shared.EphemeralContentResponse("Something went wrong!"), err
	}

	// get current page from page button
	button := actionsRow.Components[2].(*discordgo.Button)
	currentPage, _ := strconv.Atoi(string([]rune(button.Label)[0]))
	currentPage -= 1 // current page index

	// unmarshal data and update embed description
	var texts []string
	err = json.Unmarshal([]byte(existingChunks), &texts)
	if err != nil {
		rsp := shared.EphemeralContentResponse("Something went wrong!")
		return rsp, err
	}
	embed.Description = texts[currentPage-1]

	// update page button
	pageButton := NewPageButton(currentPage, len(texts), hash).Component()

	// previous page button
	if currentPage-1 > 0 {
		actionsRow.Components[0].(*discordgo.Button).Disabled = false
		actionsRow.Components[1].(*discordgo.Button).Disabled = false
	} else {
		actionsRow.Components[0].(*discordgo.Button).Disabled = true
		actionsRow.Components[1].(*discordgo.Button).Disabled = true
	}
	backwardSkipButton := actionsRow.Components[0]
	previousPageButton := actionsRow.Components[1]

	// next page button (will always be valid when using previous page button)
	actionsRow.Components[3].(*discordgo.Button).Disabled = false
	actionsRow.Components[4].(*discordgo.Button).Disabled = false
	nextPageButton := actionsRow.Components[3]
	forwardSkipButton := actionsRow.Components[4]

	// attach buttons
	buttons := []discordgo.MessageComponent{backwardSkipButton, previousPageButton, pageButton, nextPageButton, forwardSkipButton}
	rsp := shared.EphemeralEmbedResponse(embed)
	shared.AddComponent(rsp, buttons...)
	rsp.Type = discordgo.InteractionResponseUpdateMessage
	return rsp, nil
}
