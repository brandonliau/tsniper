package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"Tsniper/pkg/config"
	"Tsniper/pkg/database"
	"Tsniper/pkg/logger"

	"Tsniper/internal/command"
	"Tsniper/internal/component"
	"Tsniper/internal/manager"
	"Tsniper/internal/notifier"
	"Tsniper/internal/repository"
	"Tsniper/internal/service"

	"github.com/bwmarrin/discordgo"
	_ "modernc.org/sqlite"
)

func main() {
	// Create logger, config, and database
	logger := logger.NewStdLogger(logger.LevelDebug)
	dCfg := config.NewDiscordConfig("./config/config.yml", logger)
	sCfg := config.NewServiceConfig("./config/config.yml", logger)
	db := database.NewSqliteDB("./database.db", logger)
	defer db.Close()

	// Create new discord session
	s, err := discordgo.New("Bot " + dCfg.Token)
	if err != nil {
		logger.Fatal("Failed to create discord session : %v", err)
	}
	s.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)

	// Create repository, notifier, services, and session manager
	repo := repository.NewSnipeRepo(dCfg, sCfg, s, db, logger)
	notifier := notifier.NewDiscordNotifier(s)
	indexingService := service.NewIndexingService(sCfg, repo, db, logger)
	snipeService := service.NewSnipeService(sCfg, dCfg, s, repo, notifier, db, logger)
	m := manager.NewDiscordManager(dCfg, s, repo, logger, notifier)

	// Add event handlers
	s.AddHandler(m.CommandInteractionHandler)
	s.AddHandler(m.ComponentInteractionHandler)
	s.AddHandler(m.ReadyHandler)
	s.AddHandler(m.ResumedHandler)
	s.AddHandler(m.RateLimitHandler)
	s.AddHandler(m.GuildMemberAddHandler)
	s.AddHandler(m.GuildMemberRemoveHandler)

	// Establish websocket connection
	err = s.Open()
	if err != nil {
		logger.Fatal("Failed to establish websocket connection : %v", err)
	}
	defer s.Close()

	// Start services
	repo.Sync()
	err = indexingService.Start()
	if err != nil {
		logger.Fatal("Failed to start indexing service: %v", err)
	}
	err = snipeService.Start()
	if err != nil {
		logger.Fatal("Failed to start snipe service: %v", err)
	}

	// Register application commands
	m.RegisterCommand(command.NewAddCommand(dCfg, sCfg, repo, db))
	m.RegisterCommand(command.NewRemoveCommand(dCfg, sCfg, repo, db))
	m.RegisterCommand(command.NewClearCommand(dCfg, repo, db))
	m.RegisterCommand(command.NewCheckCommand(dCfg, repo, db))
	m.RegisterCommand(command.NewSearchCommand(dCfg, sCfg, repo, db))
	m.RegisterCommand(command.NewPingCommand())
	m.RegisterCommand(command.NewUptimeCommand(time.Now().Unix()))
	m.RegisterCommand(command.NewHelpCommand(dCfg, repo))

	// Register application components
	m.RegisterComponent(component.NewResnipeButton(dCfg, sCfg, repo, db))

	// Update bot personalization
	s.UpdateCustomStatus("👁️‍🗨️ Monitoring...")
	logger.Info("Bot running")

	// Create stop channel and block execution until a stop signal is received
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-sc

	// Stop service
	err = indexingService.Stop()
	if err != nil {
		logger.Error("Failed to stop indexing service: %v", err)
	}
	err = snipeService.Stop()
	if err != nil {
		logger.Error("Failed to stop snipe service: %v", err)
	}

	// Remove application commands
	_, err = s.ApplicationCommandBulkOverwrite(s.State.User.ID, "", nil)
	if err != nil {
		logger.Error("Failed to delete application commands")
	}
	logger.Info("Bot shut down")
}
