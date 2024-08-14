package main

import(
	"fmt"
	"time"
	
	"github.com/bwmarrin/discordgo"
)

var colorMap = map[string]int {
	"blue": 0x5865f2,
	"green": 0x2dcc70,
	"red": 0xe74d3b,
}

var emojiMap = map[string]string {
	"SPRING": ":herb:",
	"SUMMER": ":sunny:",
	"FALL": ":fallen_leaf:", 
	"WINTER": ":snowflake:",
}

// /* ADD EMBEDS */
func (runTimeData *RunTimeData) InvalidAdd(index string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title: "Invalid Request!",
		Description: fmt.Sprintf("`%s` does not exist.", index),
		Color: colorMap["red"],
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: runTimeData.Config.Image,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: time.Now().Format("01/02/2006 03:04:05 PM"),
		},
	}
}

func (runTimeData *RunTimeData) DuplicateAdd(index string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title: "Duplicate Request!",
		Description: fmt.Sprintf("You are already sniping `%s`.", index),
		Color: colorMap["red"],
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: runTimeData.Config.Image,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: time.Now().Format("01/02/2006 03:04:05 PM"),
		},
	}
}

func (runTimeData *RunTimeData) SuccessfulAdd(courseData CourseData) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title: "Successfully Added Request!",
		Description: fmt.Sprintf(
			"`%s` - %s (**Section %s**) was added to your snipe requests.",
			courseData.Index,
			courseData.Title,
			courseData.Section,
		),
		Color: colorMap["green"],
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: runTimeData.Config.Image,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: time.Now().Format("01/02/2006 03:04:05 PM"),
		},
	}
}

// /* REMOVE EMBEDS */
func (runTimeData *RunTimeData) InvalidRemove(index string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title: "Invalid Request!",
		Description: fmt.Sprintf("You are not currently sniping `%s`.", index),
		Color: colorMap["red"],
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: runTimeData.Config.Image,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: time.Now().Format("01/02/2006 03:04:05 PM"),
		},
	}
}

func (runTimeData *RunTimeData) SuccessfulRemove(courseData CourseData) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title: "Successfully Removed Request!",
		Description: fmt.Sprintf(
			"`%s` - %s (**Section %s**) was removed from your snipe requests.",
			courseData.Index,
			courseData.Title,
			courseData.Section,
		),
		Color: colorMap["green"],
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: runTimeData.Config.Image,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: time.Now().Format("01/02/2006 03:04:05 PM"),
		},
	}
}

// /* CLEAR EMBEDS */
func (runTimeData *RunTimeData) InvalidClear() *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title: "Invalid Request!",
		Description: "You have no active snipe requests.",
		Color: colorMap["red"],
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: runTimeData.Config.Image,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: time.Now().Format("01/02/2006 03:04:05 PM"),
		},
	}
}

func (runTimeData *RunTimeData) SuccessfulClear() *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title: "Success!",
		Description: "All active snipes requests have been removed.",
		Color: colorMap["green"],
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: runTimeData.Config.Image,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: time.Now().Format("01/02/2006 03:04:05 PM"),
		},
	}
}

// /* CHECK EMBEDS */
func (runTimeData *RunTimeData) InvalidCheck() *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title: "Invalid Request!",
		Description: "You have no active snipe requests.",
		Color: colorMap["red"],
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: runTimeData.Config.Image,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: time.Now().Format("01/02/2006 03:04:05 PM"),
		},
	}
}

func (runTimeData *RunTimeData) SuccessfulCheck(text string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title: "Active Requests",
		Description: text,
		Color: colorMap["blue"],
		Footer: &discordgo.MessageEmbedFooter{
			Text: time.Now().Format("01/02/2006 03:04:05 PM"),
		},
	}
}

// /* SEARCH EMBEDS */
func (runTimeData *RunTimeData) InvalidSearch(index string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title: "Invalid Request!",
		Description: fmt.Sprintf("`%s` does not exist.", index),
		Color: colorMap["red"],
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: runTimeData.Config.Image,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: time.Now().Format("01/02/2006 03:04:05 PM"),
		},
	}
}

func (runTimeData *RunTimeData) SuccessfulSearch(course CourseData) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title: fmt.Sprintf("%s (`%s`)", course.Title, course.CourseString),
		Color: colorMap["blue"],
		Fields: []*discordgo.MessageEmbedField{
			{
				Name: ":alarm_clock: Section Meeting Times",
				Value: fmt.Sprintf(">>> %s", course.Meeting),
				Inline: false,
			},
			{
				Name: "Course Name",
				Value: fmt.Sprintf("`%s`", course.Title),
				Inline: true,
			},
			{
				Name: "Index",
				Value: fmt.Sprintf("`%s`", course.Index),
				Inline: true,
			},
			{
				Name: "Section",
				Value: fmt.Sprintf("`%s`", course.Section),
				Inline: true,
			},
			{
				Name: "Instructors",
				Value: fmt.Sprintf("```fix\n%s```", course.Instructors),
				Inline: false,
			},
			{
				Name: "Special Notes",
				Value: fmt.Sprintf("```fix\n%s```", course.Notes),
				Inline: false,
			},
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: runTimeData.Config.Image,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: time.Now().Format("01/02/2006 03:04:05 PM"),
		},
	}
}

// /* DEBUG EMBEDS */
func (runTimeData *RunTimeData) HelpEmbed() *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title: "Commands",
		Color: colorMap["blue"],
		Fields: []*discordgo.MessageEmbedField{
			{
				Name: "Sniping",
				Value: fmt.Sprintf(
					"</add:%s> - Add a snipe request.\n" +
					"</remove:%s> - Remove a snipe request.\n" +
					"</clear:%s> - Remove all active snipe requests.\n" +
					"</check:%s> - View all active snipe requests.\n" +
					"</search:%s> - View course information for given index.",
					runTimeData.Registered["add"],
					runTimeData.Registered["remove"],
					runTimeData.Registered["clear"],
					runTimeData.Registered["check"],
					runTimeData.Registered["search"],
				),
				Inline: false,
			},
			{
				Name: "Miscellaneous",
				Value: fmt.Sprintf(
					"</help:%s> - List all commands.\n" +
					"</uptime:%s> - Check bot uptime.\n" +
					"</ping:%s> - Check bot latency.",
					runTimeData.Registered["help"],
					runTimeData.Registered["uptime"],
					runTimeData.Registered["ping"],
				),
				Inline: false,
			},
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: runTimeData.Config.Image,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: time.Now().Format("01/02/2006 03:04:05 PM"),
		},
	}
}

func (runTimeData *RunTimeData) UptimeEmbed() *discordgo.MessageEmbed {
	diff := time.Now().Unix() - runTimeData.StartTime
	return &discordgo.MessageEmbed{
		Title: "SnipeR Uptime",
		Description: fmt.Sprintf(
			"Last restart: <t:%d:R>\n" +
            "Uptime: %d days, %d hours, %d min, %d sec",
			runTimeData.StartTime,
			(diff / 86400),
			((diff / 3600) % 24),
			((diff / 60) % 60),
			(diff % 60),
		),
		Color: colorMap["blue"],
		Footer: &discordgo.MessageEmbedFooter{
			Text: time.Now().Format("01/02/2006 03:04:05 PM"),
		},
	}
}

// /* SNIPE EMBEDS */
func (runTimeData *RunTimeData) SnipeEmbed(course CourseData) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title: fmt.Sprintf("%s (Section %s) has opened!", course.Title, course.Section),
		Color: colorMap["blue"],
		Fields: []*discordgo.MessageEmbedField{
			{
				Name: "Course Name",
				Value: fmt.Sprintf("`%s`", course.Title),
				Inline: true,
			},
			{
				Name: "Index",
				Value: fmt.Sprintf("`%s`", course.Index),
				Inline: true,
			},
			{
				Name: "Section",
				Value: fmt.Sprintf("`%s`", course.Section),
				Inline: true,
			},
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: runTimeData.Config.Image,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: time.Now().Format("01/02/2006 03:04:05 PM"),
		},
	}
}

func (runTimeData *RunTimeData) JoinEmbed(user *discordgo.User) *discordgo.MessageEmbed {
	guild, _ := s.State.Guild(runTimeData.Config.Guild)
	return &discordgo.MessageEmbed{
		Title: "Welcome to the TSniper server!",
		Description: fmt.Sprintf("%s has joined the server!\n\nYou are user **#%d**!", user.Mention(), guild.MemberCount),
		Color: colorMap["green"],
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: user.AvatarURL(""),
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: time.Now().Format("01/02/2006 03:04:05 PM"),
		},
	}
}