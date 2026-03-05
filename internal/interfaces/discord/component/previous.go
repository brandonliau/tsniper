package component

import (
	"sort"

	"tsniper/internal/application/usecase"
	"tsniper/internal/interfaces/discord/interaction"
	"tsniper/internal/interfaces/discord/presentation"

	"tsniper/pkg/utils"

	"github.com/bwmarrin/discordgo"
)

type previousComponent struct {
	snipeService *usecase.SnipeService
}

func PreviousComponentDefinition(data ...utils.KeyValue[string, string]) discordgo.Button {
	return discordgo.Button{
		CustomID: interaction.EncodeCustomID("previous", data...),
		Label:    "⬅️",
		Style:    discordgo.PrimaryButton,
	}
}

func PreviousComponentHandler(snipeService *usecase.SnipeService) interaction.HandleFunc {
	c := &previousComponent{
		snipeService: snipeService,
	}
	return c.execute
}

func (c *previousComponent) execute(s *discordgo.Session, i *discordgo.InteractionCreate) (*discordgo.InteractionResponse, error) {
	_, params := interaction.DecodeCustomID(i.MessageComponentData().CustomID)
	cursor := params["campus"] + params["term"] + params["year"] + params["index"]

	res, err := c.snipeService.Check(usecase.CheckSnipeRequest{
		UserID: interaction.GetUserID(i),
	})
	if err != nil {
		return nil, err
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
		var embed *discordgo.MessageEmbed
		if len(res.Courses) == 0 {
			embed = presentation.InvalidCheck()
		} else {
			embed = presentation.SuccessfulCheck(res.Courses, res.Counts)
		}
		rsp := interaction.InteractionUpdateResponse(
			interaction.WithEmbeds(embed),
		)
		rsp.Data.Components = []discordgo.MessageComponent{}
		return rsp, nil
	}

	endIdx := len(res.Courses)
	for idx, crs := range res.Courses {
		key := string(crs.Scope.Campus) + string(crs.Scope.Term) + crs.Scope.Year + crs.Index
		if key >= cursor {
			endIdx = idx
			break
		}
	}

	startIdx := max(0, endIdx-presentation.MaxCoursesPerPage)
	page := res.Courses[startIdx:endIdx]

	if len(page) == 0 {
		rsp := interaction.InteractionUpdateResponse(
			interaction.WithEmbeds(presentation.InvalidCheck()),
		)
		return rsp, nil
	}

	previousBtn := PreviousComponentDefinition()
	previousBtn.Disabled = true
	if startIdx > 0 {
		first := page[0]
		previousBtn = PreviousComponentDefinition(
			utils.KeyValue[string, string]{Key: "campus", Value: string(first.Scope.Campus)},
			utils.KeyValue[string, string]{Key: "term", Value: string(first.Scope.Term)},
			utils.KeyValue[string, string]{Key: "year", Value: first.Scope.Year},
			utils.KeyValue[string, string]{Key: "index", Value: first.Index},
		)
	}

	last := page[len(page)-1]
	nextBtn := NextComponentDefinition(
		utils.KeyValue[string, string]{Key: "campus", Value: string(last.Scope.Campus)},
		utils.KeyValue[string, string]{Key: "term", Value: string(last.Scope.Term)},
		utils.KeyValue[string, string]{Key: "year", Value: last.Scope.Year},
		utils.KeyValue[string, string]{Key: "index", Value: last.Index},
	)

	if endIdx >= len(res.Courses) {
		nextBtn = NextComponentDefinition()
		nextBtn.Disabled = true
	}

	totalPages := (len(res.Courses) + presentation.MaxCoursesPerPage - 1) / presentation.MaxCoursesPerPage
	currentPage := (startIdx / presentation.MaxCoursesPerPage) + 1
	pageBtn := PageComponentDefinition(currentPage, totalPages)

	rsp := interaction.InteractionUpdateResponse(
		interaction.WithEmbeds(presentation.SuccessfulCheck(page, res.Counts)),
		interaction.WithComponents(
			previousBtn,
			pageBtn,
			nextBtn,
		),
	)
	return rsp, nil
}
