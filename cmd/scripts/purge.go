package main

import (
	"time"
	"tsniper/internal/config"
	"tsniper/internal/infrastructure/persistence/sqlite"

	"tsniper/pkg/database"
	"tsniper/pkg/logger"

	"github.com/bwmarrin/discordgo"
	_ "modernc.org/sqlite"
)

func purge() {
	logger := logger.NewStdLogger(logger.LevelInfo)

	cfg, err := config.NewYamlConfig("./config/config.yml")
	if err != nil {
		logger.Fatal("Failed to create config: %v", err)
	}

	db, err := database.NewSqliteDB("./database.db")
	if err != nil {
		logger.Fatal("Failed to create database: %v", err)
	}
	defer db.Close()

	err = sqlite.Migrate(db)
	if err != nil {
		logger.Fatal("Failed to run migrations: %v", err)
	}

	s, err := discordgo.New("Bot " + cfg.Discord.Token)
	if err != nil {
		logger.Fatal("Failed to create discord session: %v", err)
	}
	s.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)

	err = s.Open()
	if err != nil {
		logger.Fatal("Failed to establish websocket connection: %v", err)
	}
	defer s.Close()

	botID := s.State.User.ID

	userRepository := sqlite.NewUserRepository(db)
	users, err := userRepository.GetAll()
	if err != nil {
		logger.Fatal("Failed to get all users: %v", err)
	}

	totalDeleted := 0
	for _, usr := range users {
		ch, err := s.UserChannelCreate(usr.ID)
		if err != nil {
			logger.Warn("Failed to open DM channel for user %s: %v", usr.ID, err)
			continue
		}

		deleted := 0
		var beforeID string
		for {
			messages, err := s.ChannelMessages(ch.ID, 100, beforeID, "", "")
			if err != nil {
				logger.Error("Failed to fetch messages for user %s: %v", usr.ID, err)
				break
			}
			if len(messages) == 0 {
				break
			}

			for _, msg := range messages {
				if msg.Author.ID == botID {
					err := s.ChannelMessageDelete(ch.ID, msg.ID)
					if err != nil {
						logger.Error("Failed to delete message %s for user %s: %v", msg.ID, usr.ID, err)
					} else {
						deleted++
					}
					time.Sleep(500 * time.Millisecond)
				}
			}

			beforeID = messages[len(messages)-1].ID
			if len(messages) < 100 {
				break
			}
		}

		if deleted > 0 {
			logger.Info("Deleted %d messages for user %s", deleted, usr.ID)
			totalDeleted += deleted
		}
	}

	logger.Info("Purged %d DM messages total", totalDeleted)
}
