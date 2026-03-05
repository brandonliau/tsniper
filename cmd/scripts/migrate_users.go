package main

import (
	"tsniper/internal/config"
	"tsniper/internal/domain/scope"
	"tsniper/internal/domain/user"
	"tsniper/internal/infrastructure/persistence/sqlite"

	"tsniper/pkg/database"
	"tsniper/pkg/logger"

	"github.com/bwmarrin/discordgo"
	_ "modernc.org/sqlite"
)

func migrateUsers() {
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

	userRepository := sqlite.NewUserRepository(db)

	var lastUserID string
	updated := 0
	for {
		members, _ := s.GuildMembers(cfg.Discord.GuildID, lastUserID, 1000)
		if len(members) < 1 {
			break
		}

		for _, member := range members {
			usr := user.NewUser(member.User.ID)
			if err := userRepository.Create(usr); err != nil {
				logger.Warn("Failed to create user %s: %v", member.User.ID, err)
				continue
			}

			for _, roleID := range member.Roles {
				roleName, ok := cfg.Discord.Roles[roleID]
				if !ok {
					continue
				}

				campus, err := scope.ParseCampus(roleName)
				if err != nil {
					logger.Warn("Invalid campus role %s for user %s: %v", roleName, member.User.ID, err)
					continue
				}

				usr.SetCampus(campus)
				if err := userRepository.Save(usr); err != nil {
					logger.Warn("Failed to set campus %s for user %s: %v", roleName, member.User.ID, err)
					continue
				}

				updated++
				logger.Info("Set campus %s for user %s", roleName, member.User.ID)
			}
		}

		lastUserID = members[len(members)-1].User.ID
		if len(members) < 1000 {
			break
		}
	}

	logger.Info("Migration complete: updated %d users", updated)
}
