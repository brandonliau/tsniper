package component

import (
	"encoding/json"
	"strconv"

	"Tsniper/internal/repository"
	"Tsniper/internal/shared"

	"Tsniper/pkg/database"

	"github.com/bwmarrin/discordgo"
)

type nextPageButton struct {
	disabled bool
	repo     repository.Repository
	db       database.Database
}

func NewNextPageButton(disabled bool, repo repository.Repository, db database.Database) *nextPageButton {
	return &nextPageButton{
		disabled: disabled,
		repo:     repo,
		db:       db,
	}
}

func (c *nextPageButton) CustomID() string {
	return c.Component().(discordgo.Button).CustomID
}

func (c *nextPageButton) Component() discordgo.MessageComponent {
	button := discordgo.Button{
		Emoji: &discordgo.ComponentEmoji{
			Name: "➡️",
		},
		Style:    discordgo.PrimaryButton,
		CustomID: "nextPage",
	}
	button.Disabled = c.disabled
	return button
}

func (c *nextPageButton) Execute(args *shared.CmdArgs) (*discordgo.InteractionResponse, error) {
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
	embed.Description = texts[currentPage+1]

	// update page button
	pageButton := NewPageButton(currentPage+2, len(texts), hash).Component()

	// previous buttons (will always be valid when using the forward skip button)
	actionsRow.Components[0].(*discordgo.Button).Disabled = false
	actionsRow.Components[1].(*discordgo.Button).Disabled = false
	backwardSkipButton := actionsRow.Components[0]
	previousPageButton := actionsRow.Components[1]

	// next page button
	if currentPage+2 < len(texts) {
		actionsRow.Components[3].(*discordgo.Button).Disabled = false
		actionsRow.Components[4].(*discordgo.Button).Disabled = false
	} else {
		actionsRow.Components[3].(*discordgo.Button).Disabled = true
		actionsRow.Components[4].(*discordgo.Button).Disabled = true
	}
	nextPageButton := actionsRow.Components[3]
	forwardSkipButton := actionsRow.Components[4]

	// attach buttons
	buttons := []discordgo.MessageComponent{backwardSkipButton, previousPageButton, pageButton, nextPageButton, forwardSkipButton}
	rsp := shared.EphemeralEmbedResponse(embed)
	shared.AddComponent(rsp, buttons...)
	rsp.Type = discordgo.InteractionResponseUpdateMessage
	return rsp, nil
}
