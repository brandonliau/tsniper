package presentation

import (
	"fmt"
	"strings"
	"time"

	"tsniper/internal/application/view"
	"tsniper/internal/domain/scope"

	"github.com/bwmarrin/discordgo"
)

func SuccessfulCheck(courses []*view.CourseView, counts []int) *discordgo.MessageEmbed {
	var builder strings.Builder
	for i, crs := range courses {
		fmt.Fprintf(
			&builder,
			"%s `%s` - %s (**Section %s**) | :eyes: `%d` | %s",
			EmojiMap[scope.Term(crs.Term).DisplayName()],
			crs.Index,
			crs.Title,
			crs.Section,
			counts[i],
			LastOpenDisplayString(crs.LastOpen),
		)
	}
	return &discordgo.MessageEmbed{
		Title:       "Active Requests",
		Description: builder.String(),
		Color:       Blue,
		Footer: &discordgo.MessageEmbedFooter{
			Text: time.Now().Format("01/02/2006 03:04:05 PM"),
		},
	}
}

func InvalidCheck() *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "Invalid Request!",
		Description: "You have no active snipe requests.",
		Color:       Red,
		Footer: &discordgo.MessageEmbedFooter{
			Text: time.Now().Format("01/02/2006 03:04:05 PM"),
		},
	}
}
