package main

import (
	"tsniper/internal/application/usecase"
	"tsniper/internal/config"
	"tsniper/internal/infrastructure/persistence/sqlite"

	"tsniper/pkg/database"
	"tsniper/pkg/eventbus"
	"tsniper/pkg/logger"
	"tsniper/pkg/utils"

	"tsniper/internal/application/event"
	"tsniper/internal/infrastructure/persistence/memory"

	"github.com/bwmarrin/discordgo"
	_ "modernc.org/sqlite"
)

func clean() {
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

	err = sqlite.Migrate(db)
	if err != nil {
		logger.Fatal("Failed to run migrations: %v", err)
	}
	err = db.Exec("DROP TABLE IF EXISTS snipes")
	if err != nil {
		logger.Fatal("Failed to drop snipes table: %v", err)
	}
	err = db.Exec(`CREATE TABLE IF NOT EXISTS snipes (
		user_id      TEXT,
		course_index TEXT,
		campus       TEXT,
		term         TEXT,
		year         TEXT,
		created_at   BIGINT,
		pending      BOOLEAN,
		UNIQUE(user_id, course_index, campus, term, year)
	)`)
	if err != nil {
		logger.Fatal("Failed to recreate snipes table: %v", err)
	}
	err = db.Exec("CREATE INDEX IF NOT EXISTS idx_snipes_user_id ON snipes(user_id)")
	if err != nil {
		logger.Fatal("Failed to create snipes user_id index: %v", err)
	}
	err = db.Exec("CREATE INDEX IF NOT EXISTS idx_snipes_course_index ON snipes(course_index)")
	if err != nil {
		logger.Fatal("Failed to create snipes course_index index: %v", err)
	}

	entitlementUpdateEventBus := eventbus.NewEventBus[event.EntitlementUpdate](50)
	snipeCache := memory.NewSnipeCache()
	snipeRepository := sqlite.NewSnipeRepository(db)
	userRepository := sqlite.NewUserRepository(db)
	userService := usecase.NewUserService(entitlementUpdateEventBus, snipeCache, snipeRepository, userRepository)

	res, err := userService.GetAll(usecase.GetAllUsersRequest{})
	if err != nil {
		logger.Fatal("Failed to get all users: %v", err)
	}

	for _, usr := range res.Users {
		if usr.Entitlement.IsUnlimited() || !usr.Entitlement.IsActive() {
			continue
		}

		_, err := userService.RevokeSemesters(usecase.RevokeSemestersRequest{
			UserID: usr.ID,
			Count:  1,
		})
		if err != nil {
			logger.Error("Failed to decrement semesters for user %s: %v", usr.ID, err)
		}
	}

	// Sync guild members with database
	var members []string
	var lastUserID string
	for {
		gm, _ := s.GuildMembers(cfg.Discord.GuildID, lastUserID, 1000)
		if len(gm) < 1 {
			break
		}
		for _, member := range gm {
			members = append(members, member.User.ID)
		}
		lastUserID = gm[len(gm)-1].User.ID
	}

	var users []string
	for _, usr := range res.Users {
		users = append(users, usr.ID)
	}

	toAdd := utils.Difference(members, users)
	toRemove := utils.Difference(users, members)

	for _, userID := range toAdd {
		_, err := userService.Join(usecase.UserJoinRequest{UserID: userID})
		if err != nil {
			logger.Warn("Failed to add user %s: %v", userID, err)
		}
	}

	for _, userID := range toRemove {
		_, err := userService.Leave(usecase.UserLeaveRequest{UserID: userID})
		if err != nil {
			logger.Warn("Failed to remove user %s: %v", userID, err)
		}
	}

	allUsers, err := userService.GetAll(usecase.GetAllUsersRequest{})
	if err != nil {
		logger.Fatal("Failed to get all users for role sync: %v", err)
	}

	for _, usr := range allUsers.Users {
		switch {
		case usr.Entitlement.IsUnlimited():
			if err := s.GuildMemberRoleAdd(cfg.Discord.GuildID, usr.ID, cfg.Discord.Roles.Unlimited); err != nil {
				logger.Error("Failed to add unlimited role for user %s: %v", usr.ID, err)
			}
			if err := s.GuildMemberRoleRemove(cfg.Discord.GuildID, usr.ID, cfg.Discord.Roles.Standard); err != nil {
				logger.Error("Failed to remove standard role for unlimited user %s: %v", usr.ID, err)
			}
		case usr.Entitlement.IsActive():
			if err := s.GuildMemberRoleAdd(cfg.Discord.GuildID, usr.ID, cfg.Discord.Roles.Standard); err != nil {
				logger.Error("Failed to add standard role for user %s: %v", usr.ID, err)
			}
			if err := s.GuildMemberRoleRemove(cfg.Discord.GuildID, usr.ID, cfg.Discord.Roles.Unlimited); err != nil {
				logger.Error("Failed to remove unlimited role for standard user %s: %v", usr.ID, err)
			}
		default:
			if err := s.GuildMemberRoleRemove(cfg.Discord.GuildID, usr.ID, cfg.Discord.Roles.Standard); err != nil {
				logger.Error("Failed to remove standard role for inactive user %s: %v", usr.ID, err)
			}
			if err := s.GuildMemberRoleRemove(cfg.Discord.GuildID, usr.ID, cfg.Discord.Roles.Unlimited); err != nil {
				logger.Error("Failed to remove unlimited role for inactive user %s: %v", usr.ID, err)
			}
		}
	}

	entitlementUpdateEventBus.Close()
	logger.Info("Database cleaned")
}
