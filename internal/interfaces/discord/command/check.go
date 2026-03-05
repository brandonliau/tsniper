package command

import (
	"sort"

	"tsniper/internal/application/usecase"
	"tsniper/internal/interfaces/discord/component"
	"tsniper/internal/interfaces/discord/interaction"
	"tsniper/internal/interfaces/discord/presentation"

	"tsniper/pkg/utils"

	"github.com/bwmarrin/discordgo"
)

type checkCommand struct {
	snipeService *usecase.SnipeService
}

func CheckCommandDefinition() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "check",
		Description: "View all active snipe requests.",
	}
}

func CheckCommandHandler(snipeService *usecase.SnipeService) interaction.HandleFunc {
	c := &checkCommand{
		snipeService: snipeService,
	}
	return c.execute
}

func (c *checkCommand) execute(s *discordgo.Session, i *discordgo.InteractionCreate) (*discordgo.InteractionResponse, error) {
	req := usecase.CheckSnipeRequest{
		UserID: interaction.GetUserID(i),
	}

	res, err := c.snipeService.Check(req)
	if err != nil {
		return nil, err
	}

	if len(res.Courses) == 0 {
		rsp := interaction.InteractionInitialResponse(
			interaction.WithEmbeds(presentation.InvalidCheck()),
			interaction.WithEphemeral(),
		)
		return rsp, nil
	}

	sort.SliceStable(res.Courses, func(i, j int) bool {
		a, b := res.Courses[i], res.Courses[j]
		if a.Scope.Campus != b.Scope.Campus {
			return a.Scope.Campus < b.Scope.Campus
		}
		if a.Scope.Term != b.Scope.Term {
			return a.Scope.Term < b.Scope.Term
		}
		if a.Scope.Year != b.Scope.Year {
			return a.Scope.Year < b.Scope.Year
		}
		return a.Index < b.Index
	})

	if len(res.Courses) <= presentation.MaxCoursesPerPage {
		rsp := interaction.InteractionInitialResponse(
			interaction.WithEmbeds(presentation.SuccessfulCheck(res.Courses, res.Counts)),
			interaction.WithEphemeral(),
		)
		return rsp, nil
	}

	endIndex := min(presentation.MaxCoursesPerPage, len(res.Courses))
	page := res.Courses[:endIndex]

	last := page[len(page)-1]
	nextBtn := component.NextComponentDefinition(
		utils.KeyValue[string, string]{Key: "campus", Value: string(last.Scope.Campus)},
		utils.KeyValue[string, string]{Key: "term", Value: string(last.Scope.Term)},
		utils.KeyValue[string, string]{Key: "year", Value: last.Scope.Year},
		utils.KeyValue[string, string]{Key: "index", Value: last.Index},
	)

	previousBtn := component.PreviousComponentDefinition()
	previousBtn.Disabled = true

	totalPages := (len(res.Courses) + presentation.MaxCoursesPerPage - 1) / presentation.MaxCoursesPerPage
	pageBtn := component.PageComponentDefinition(1, totalPages)

	rsp := interaction.InteractionInitialResponse(
		interaction.WithEmbeds(presentation.SuccessfulCheck(page, res.Counts)),
		interaction.WithComponents(
			previousBtn,
			pageBtn,
			nextBtn,
		),
		interaction.WithEphemeral(),
	)
	return rsp, nil
}
