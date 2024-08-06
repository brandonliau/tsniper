package main

import (
	"strings"
	"github.com/bwmarrin/discordgo"
)

var Commands []*discordgo.ApplicationCommand

func (runTimeData RunTimeData) SeasonChoiceOptions() []*discordgo.ApplicationCommandOptionChoice {
	seasonChoices := make([]*discordgo.ApplicationCommandOptionChoice, 0)
	for _, season := range runTimeData.Config.CurrentSeasons {
		data := discordgo.ApplicationCommandOptionChoice{
			Name: strings.ToLower(season),
			Value: season,
		}
		seasonChoices = append(seasonChoices, &data)
	}
	return seasonChoices
}

func (runTimeData RunTimeData) CampusChoiceOptions() []*discordgo.ApplicationCommandOptionChoice {
	campusName := map[string]string{
		"NB": "new brunswick",
		"NK": "newark",
		"CM": "camden",
	}
	campusChoices := make([]*discordgo.ApplicationCommandOptionChoice, 0)
	for _, campus := range runTimeData.Config.CurrentCampuses {
		data := discordgo.ApplicationCommandOptionChoice{
			Name: campusName[campus],
			Value: campus,
		}
		campusChoices = append(campusChoices, &data)
	}
	return campusChoices
}

func (runTimeData RunTimeData) InitCommands() {
	Commands = []*discordgo.ApplicationCommand{
		{
			Name: "add",
			Description: "Add a snipe request.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type: discordgo.ApplicationCommandOptionString,
					Name: "index",
					Description: "index",
					Required: true,
					MaxLength: 5,
				},
				{
					Type: discordgo.ApplicationCommandOptionString,
					Name: "season",
					Description: "season",
					Required: false,
					Choices: runTimeData.SeasonChoiceOptions(),
				},
				{
					Type: discordgo.ApplicationCommandOptionString,
					Name: "campus",
					Description: "campus",
					Required: false,
					Choices: runTimeData.CampusChoiceOptions(),
				},
			},
		},
		{
			Name: "remove",
			Description: "Remove a snipe request.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type: discordgo.ApplicationCommandOptionString,
					Name: "index",
					Description: "index",
					Required: true,
					MaxLength: 5,
				},
				{
					Type: discordgo.ApplicationCommandOptionString,
					Name: "season",
					Description: "season",
					Required: false,
					Choices: runTimeData.SeasonChoiceOptions(),
				},
				{
					Type: discordgo.ApplicationCommandOptionString,
					Name: "campus",
					Description: "campus",
					Required: false,
					Choices: runTimeData.CampusChoiceOptions(),
				},
			},
		},
		{
			Name: "clear",
			Description: "Remove all active snipe requests.",
		},
		{
			Name: "check",
			Description: "View all active snipe requests.",
		},
		{
			Name: "search",
			Description: "View course information for given index.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type: discordgo.ApplicationCommandOptionString,
					Name: "index",
					Description: "index",
					Required: true,
					MaxLength: 5,
				},
				{
					Type: discordgo.ApplicationCommandOptionString,
					Name: "season",
					Description: "season",
					Required: false,
					Choices: runTimeData.SeasonChoiceOptions(),
				},
				{
					Type: discordgo.ApplicationCommandOptionString,
					Name: "campus",
					Description: "campus",
					Required: false,
					Choices: runTimeData.CampusChoiceOptions(),
				},
			},
		},
		{
			Name: "help",
			Description: "List all commands.",
		},
		{
			Name: "uptime",
			Description: "Check bot uptime.",
		},
		{
			Name: "ping",
			Description: "Check bot latency.",
		},
	}
}
