package main

import (
	"fmt"
	"net/url"
	
	"github.com/bwmarrin/discordgo"
)

var CommandHandlers map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)
var ComponentHandlers map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)

func (runTimeData *RunTimeData) InitHandlers() {
	CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"add": func(s *discordgo.Session, i *discordgo.InteractionCreate) {runTimeData.AddCommandHandler(s, i)},
		"remove": func(s *discordgo.Session, i *discordgo.InteractionCreate) {runTimeData.RemoveCommandHandler(s, i)},
		"clear": func(s *discordgo.Session, i *discordgo.InteractionCreate) {runTimeData.ClearCommandHandler(s, i)},
		"check": func(s *discordgo.Session, i *discordgo.InteractionCreate) {runTimeData.CheckCommandHandler(s, i)},
		"search": func(s *discordgo.Session, i *discordgo.InteractionCreate) {runTimeData.SearchCommandHandler(s, i)},
		"help": func(s *discordgo.Session, i *discordgo.InteractionCreate) {runTimeData.HelpCommandHandler(s, i)},
		"uptime": func(s *discordgo.Session, i *discordgo.InteractionCreate) {runTimeData.UptimeCommandHandler(s, i)},
		"ping": func(s *discordgo.Session, i *discordgo.InteractionCreate) {runTimeData.PingCommandHandler(s, i)},
	}
	ComponentHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"resnipe": func(s *discordgo.Session, i *discordgo.InteractionCreate) { runTimeData.ResnipeButtonHandler(s, i)},
	}
}

func (runTimeData *RunTimeData) AddCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var memberID, index string
	var embed *discordgo.MessageEmbed
	if i.Member != nil {
		memberID = i.Member.User.ID
	} else {
		memberID = i.User.ID
	}
	campus := runTimeData.GetCampus(memberID)
	season := runTimeData.Config.CurrentSeasons[0]
	for _, options := range i.ApplicationCommandData().Options {
		switch options.Name {
		case "index":
			index = options.Value.(string)
		case "season":
			season = options.Value.(string)
		case "campus":
			campus = options.Value.(string)
		}
	}
	if ok := runTimeData.CheckCourseExist(index, campus, season); !ok {
		embed = runTimeData.InvalidAdd(index)
	} else if runTimeData.AreadySniping(memberID, index, campus, season) {
		embed = runTimeData.DuplicateAdd(index)
	} else {
		runTimeData.AddSnipe(index, memberID, campus, season)
		runTimeData.UpdateTracking(1, index, campus, season)
		course := runTimeData.GetCourseData(index, campus, season)
		embed = runTimeData.SuccessfulAdd(course)
	}
	SendEmbedReponse(s, i, embed)
}

func (runTimeData *RunTimeData) RemoveCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var memberID, index string
	var embed *discordgo.MessageEmbed
	if i.Member != nil {
		memberID = i.Member.User.ID
	} else {
		memberID = i.User.ID
	}
	campus := runTimeData.GetCampus(memberID)
	season := runTimeData.Config.CurrentSeasons[0]
	for _, options := range i.ApplicationCommandData().Options {
		switch options.Name {
		case "index":
			index = options.Value.(string)
		case "season":
			season = options.Value.(string)
		case "campus":
			campus = options.Value.(string)
		}
	}
	if !runTimeData.AreadySniping(memberID, index, campus, season) {
		embed = runTimeData.InvalidRemove(index)
	} else {
		runTimeData.RemoveSnipe(index, memberID, season)
		runTimeData.UpdateTracking(2, index, campus, season)
		course := runTimeData.GetCourseData(index, campus, season)
		embed = runTimeData.SuccessfulRemove(course)
	}
	SendEmbedReponse(s, i, embed)
}

func (runTimeData *RunTimeData) ClearCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var memberID string
	var embed *discordgo.MessageEmbed
	if i.Member != nil {
		memberID = i.Member.User.ID
	} else {
		memberID = i.User.ID
	}
	snipes := runTimeData.GetSnipes(memberID)
	if len(snipes) == 0 {
		embed = runTimeData.InvalidClear()
	} else {
		runTimeData.ClearSnipe(memberID)
		runTimeData.SyncTracking()
		embed = runTimeData.SuccessfulClear()
	}
	SendEmbedReponse(s, i, embed)
}

func (runTimeData *RunTimeData) CheckCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var memberID string
	var embed *discordgo.MessageEmbed
	if i.Member != nil {
		memberID = i.Member.User.ID
	} else {
		memberID = i.User.ID
	}
	snipes := runTimeData.GetSnipes(memberID)
	if len(snipes) != 0 {
		text := ""
		for _, course := range snipes {
			index := course[0]
			campus := course[1]
			season := course[2]
			data := runTimeData.GetCourseData(index, campus, season)
			count := runTimeData.GetSnipeCount(index, campus, season)
			lastOpen := runTimeData.GetLastOpen(data.Index, campus, season)
			if lastOpen == -1 {
				text += fmt.Sprintf(
					"%s `%s` - %s (**Section %s**) | :eyes: `%d` | `Unknown`\n",
					emojiMap[season],
					data.Index,
					data.Title,
					data.Section,
					count,
				)
			} else {
				text += fmt.Sprintf(
					"%s `%s` - %s (**Section %s**) | :eyes: `%d` | <t:%d:R>\n",
					emojiMap[season],
					data.Index,
					data.Title,
					data.Section,
					count,
					lastOpen,
				)
			}
		}
		embed = runTimeData.SuccessfulCheck(text)
	} else {
		embed = runTimeData.InvalidCheck()
	}
	SendEmbedReponse(s, i, embed)
}

func (runTimeData *RunTimeData) SearchCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var memberID, index string
	var embed *discordgo.MessageEmbed
	if i.Member != nil {
		memberID = i.Member.User.ID
	} else {
		memberID = i.User.ID
	}
	campus := runTimeData.GetCampus(memberID)
	season := runTimeData.Config.CurrentSeasons[0]
	for _, options := range i.ApplicationCommandData().Options {
		switch options.Name {
		case "index":
			index = options.Value.(string)
		case "season":
			season = options.Value.(string)
		case "campus":
			campus = options.Value.(string)
		}
	}
	if ok := runTimeData.CheckCourseExist(index, campus, season); ok {
		course := runTimeData.GetCourseData(index, campus, season)
		embed = runTimeData.SuccessfulSearch(course)
	} else {
		embed = runTimeData.InvalidSearch(index)
	}
	SendEmbedReponse(s, i, embed)
}

func (runTimeData *RunTimeData) HelpCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	SendEmbedReponse(s, i, runTimeData.HelpEmbed())
}

func (runTimeData *RunTimeData) UptimeCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	SendEmbedReponse(s, i, runTimeData.UptimeEmbed())
}

func (runTimeData *RunTimeData) PingCommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	SendMessageResponse(s, i, fmt.Sprintf("Pong! `%d ms`", int(s.HeartbeatLatency().Milliseconds())))
}

func (runTimeData *RunTimeData) ResnipeButtonHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	memberID := i.Interaction.User.ID
	components := i.Message.Components[0]
	actionsRow := components.(*discordgo.ActionsRow)
	button := actionsRow.Components[0].(*discordgo.Button)
	resnipeURL, _ := url.Parse(button.URL)
	index := resnipeURL.Query().Get("indexList")
	campus := resnipeURL.Query().Get("campus")
	season := runTimeData.GetSeason(resnipeURL.Query().Get("semesterSelection"))
	if runTimeData.AreadySniping(memberID, index, campus, season) {
		SendMessageResponse(s, i, fmt.Sprintf("You are already sniping `%s`.", index))
	} else {
		runTimeData.AddSnipe(index, memberID, campus, season)
		runTimeData.UpdateTracking(1, index, campus, season)
		SendMessageResponse(s, i, fmt.Sprintf("Successfully re-added `%s` to your snipe requests.", index))
	}
}
