package component

import (
	"tsniper/internal/application/usecase"
	"tsniper/internal/application/view"
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

	view.SortCourseViews(res.Courses, res.Counts)

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
		key := crs.Campus + crs.Term + crs.Year + crs.Index
		if key >= cursor {
			endIdx = idx
			break
		}
	}

	startIdx := max(0, endIdx-presentation.MaxCoursesPerPage)
	pageCourses := res.Courses[startIdx:endIdx]
	pageCounts := res.Counts[startIdx:endIdx]

	if len(pageCourses) == 0 {
		rsp := interaction.InteractionUpdateResponse(
			interaction.WithEmbeds(presentation.InvalidCheck()),
		)
		return rsp, nil
	}

	previousBtn := PreviousComponentDefinition()
	previousBtn.Disabled = true
	if startIdx > 0 {
		first := pageCourses[0]
		previousBtn = PreviousComponentDefinition(
			utils.KeyValue[string, string]{Key: "campus", Value: first.Campus},
			utils.KeyValue[string, string]{Key: "term", Value: first.Term},
			utils.KeyValue[string, string]{Key: "year", Value: first.Year},
			utils.KeyValue[string, string]{Key: "index", Value: first.Index},
		)
	}

	last := pageCourses[len(pageCourses)-1]
	nextBtn := NextComponentDefinition(
		utils.KeyValue[string, string]{Key: "campus", Value: last.Campus},
		utils.KeyValue[string, string]{Key: "term", Value: last.Term},
		utils.KeyValue[string, string]{Key: "year", Value: last.Year},
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
		interaction.WithEmbeds(presentation.SuccessfulCheck(pageCourses, pageCounts)),
		interaction.WithComponents(
			previousBtn,
			pageBtn,
			nextBtn,
		),
	)
	return rsp, nil
}
