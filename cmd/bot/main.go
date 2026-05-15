package main

import (
	"os"
	"os/signal"
	"syscall"

	"tsniper/internal/application/event"
	"tsniper/internal/application/usecase"
	"tsniper/internal/application/view"
	"tsniper/internal/application/worker"
	"tsniper/internal/config"
	"tsniper/internal/domain/scope"
	"tsniper/internal/infrastructure/external"
	"tsniper/internal/infrastructure/persistence/memory"
	"tsniper/internal/infrastructure/persistence/sqlite"
	"tsniper/internal/infrastructure/rutgers"
	"tsniper/internal/interfaces/discord"
	"tsniper/internal/interfaces/discord/command"
	"tsniper/internal/interfaces/discord/component"

	"tsniper/pkg/database"
	"tsniper/pkg/eventbus"
	"tsniper/pkg/logger"

	"github.com/bwmarrin/discordgo"
	_ "modernc.org/sqlite"
)

func main() {
	// Create logger
	logger := logger.NewStdLogger(logger.LevelInfo)

	// Create config
	cfg, err := config.NewYamlConfig("./config/config.yml")
	if err != nil {
		logger.Fatal("Failed to create config: %v", err)
	}

	// Create database
	db, err := database.NewSqliteDB("./database.db")
	if err != nil {
		logger.Fatal("Failed to create database: %v", err)
	}
	defer db.Close()

	// Perform database migrations
	err = sqlite.Migrate(db)
	if err != nil {
		logger.Fatal("Failed to migrate database: %v", err)
	}

	// Create discord session
	s, err := discordgo.New("Bot " + cfg.Discord.Token)
	if err != nil {
		logger.Fatal("Failed to create discord session: %v", err)
	}
	s.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)

	// Create event buses
	openCourseEventBus := eventbus.NewEventBus[event.CourseOpen](100)

	// Create domain models
	activeScope := scope.NewActiveScope(cfg.Academic.Campuses, cfg.Academic.Seasons)
	seasons := view.FromSeasons(activeScope.Seasons())
	campuses := view.FromCampuses(activeScope.Campuses())

	// Create infrastructure repositories
	snipeRepository := sqlite.NewSnipeRepository(db)
	courseRepository := sqlite.NewCourseRepository(db)
	userRepository := sqlite.NewUserRepository(db)

	// Create infrastructure caches
	snipeCache := memory.NewSnipeCache()

	// Create ports
	systemMonitor := external.NewSystemMonitor()
	courseFeed := rutgers.NewCourseFeed()
	sectionsFeed := rutgers.NewSectionsFeed()

	// Create application usescases
	snipeService := usecase.NewSnipeService(activeScope, snipeCache, snipeRepository, courseRepository, userRepository)
	courseService := usecase.NewCourseService(activeScope, snipeRepository, userRepository, courseRepository)
	userService := usecase.NewUserService(snipeCache, snipeRepository, userRepository)
	systemService := usecase.NewSystemService(systemMonitor)

	// Create application gateways
	discordGateway := discord.NewGateway(s, cfg.Discord.ApplicationID, cfg.Discord.GuildID, userService, cfg.Discord, logger)

	// Create event notifiers
	openCourseNotifier := discord.NewOpenCourseNotifier(openCourseEventBus, s, cfg.Customization, logger)

	// Register application commands
	discordGateway.RegisterCommand(command.AddCommandDefinition(seasons, campuses), command.AddCommandHandler(snipeService, cfg.Customization))
	discordGateway.RegisterCommand(command.RemoveCommandDefinition(seasons, campuses), command.RemoveCommandHandler(snipeService, cfg.Customization))
	discordGateway.RegisterCommand(command.ClearCommandDefinition(), command.ClearCommandHandler(snipeService, cfg.Customization))
	discordGateway.RegisterCommand(command.CheckCommandDefinition(), command.CheckCommandHandler(snipeService))
	discordGateway.RegisterCommand(command.SearchCommandDefinition(seasons, campuses), command.SearchCommandHandler(courseService, cfg.Customization))
	discordGateway.RegisterCommand(command.StatusCommandDefinition(), command.StatusCommandHandler(systemService, cfg.Customization))
	discordGateway.RegisterCommand(command.HelpCommandDefinition(), command.HelpCommandHandler(cfg.Discord.ApplicationID, cfg.Customization))

	// Register application components
	discordGateway.RegisterComponent(component.ResnipeComponentDefinition(), component.ResnipeComponentHandler(snipeService))
	discordGateway.RegisterComponent(component.NextComponentDefinition(), component.NextComponentHandler(snipeService))
	discordGateway.RegisterComponent(component.PreviousComponentDefinition(), component.PreviousComponentHandler(snipeService))

	// Register application services
	orchestrator := worker.NewOrchestrator(logger)
	orchestrator.RegisterWorker("course_indexer", worker.NewCourseIndexer(activeScope, courseFeed, courseRepository, logger))
	orchestrator.RegisterWorker("snipe_monitor", worker.NewMonitorService(activeScope, openCourseEventBus, sectionsFeed, snipeCache, snipeRepository, courseRepository, userRepository, logger))

	// Add event handlers
	s.AddHandler(discordGateway.InteractionHandler)
	s.AddHandler(discordGateway.ReadyHandler)
	s.AddHandler(discordGateway.ResumedHandler)
	s.AddHandler(discordGateway.RateLimitHandler)
	s.AddHandler(discordGateway.MemberJoinHandler)
	s.AddHandler(discordGateway.MemberLeaveHandler)
	s.AddHandler(discordGateway.MemberUpdateHandler)

	// Start gateway
	discordGateway.Start()

	// Start notifiers
	openCourseNotifier.Start()

	// Start workers
	err = orchestrator.StartAll()
	if err != nil {
		logger.Fatal("Failed to start services: %v", err)
	}

	// Establish websocket connection
	err = s.Open()
	if err != nil {
		logger.Fatal("Failed to establish websocket connection : %v", err)
	}
	defer s.Close()

	// Bot online
	logger.Info("Bot running")

	// Create stop channel and block execution until a stop signal is received
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-sc

	// Stop workers
	err = orchestrator.StopAll()
	if err != nil {
		logger.Fatal("Failed to stop services: %v", err)
	}

	// Stop application gateways
	discordGateway.Stop()

	// Close event buses
	openCourseEventBus.Close()

	// Stop event notifiers
	openCourseNotifier.Stop()

	// Remove application commands
	_, err = s.ApplicationCommandBulkOverwrite(s.State.User.ID, "", nil)
	if err != nil {
		logger.Error("Failed to delete application commands")
	}
	logger.Info("Bot shut down")
}
