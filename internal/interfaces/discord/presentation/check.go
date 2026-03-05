package presentation

import (
	"fmt"
	"strings"
	"time"

	"tsniper/internal/domain/course"

	"github.com/bwmarrin/discordgo"
)

func SuccessfulCheck(courses []*course.Course, counts map[*course.Course]int) *discordgo.MessageEmbed {
	var builder strings.Builder
	for _, crs := range courses {
		fmt.Fprintf(
			&builder,
			"%s `%s` - %s (**Section %s**) | :eyes: `%d` | %s",
			EmojiMap[crs.Scope.Term.DisplayName()],
			crs.Index,
			crs.Title,
			crs.Section,
			counts[crs],
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
