package discord

import (
	"strings"
	"tsniper/internal/application/usecase"
	"tsniper/internal/config"
	"tsniper/internal/interfaces/discord/interaction"

	"tsniper/pkg/logger"

	"github.com/bwmarrin/discordgo"
)

type gateway struct {
	session       *discordgo.Session
	applicationID string
	guildID       string
	userService   *usecase.UserService
	handleFuncs   map[string]interaction.HandleFunc
	cfg           *config.DiscordConfig
	logger        logger.Logger
}

func NewGateway(session *discordgo.Session, applicationID string, guildID string, userService *usecase.UserService, cfg *config.DiscordConfig, logger logger.Logger) *gateway {
	return &gateway{
		session:       session,
		applicationID: applicationID,
		guildID:       guildID,
		userService:   userService,
		handleFuncs:   make(map[string]interaction.HandleFunc),
		cfg:           cfg,
		logger:        logger,
	}
}

func (g *gateway) Start() {
	res, err := g.userService.GetAll(usecase.GetAllUsersRequest{})
	if err != nil {
		g.logger.Fatal("Failed to start discord gateway: %v", err)
	}

	userSet := make(map[string]struct{}, len(res.Users))
	for _, usr := range res.Users {
		userSet[usr.ID] = struct{}{}
	}

	roleToCampus := map[string]string{
		g.cfg.Roles.NB: "NB",
		g.cfg.Roles.NK: "NK",
		g.cfg.Roles.CM: "CM",
	}

	memberRoles := make(map[string][]string)
	var lastUserID string
	for {
		gm, _ := g.session.GuildMembers(g.guildID, lastUserID, 1000)
		if len(gm) < 1 {
			break
		}
		for _, member := range gm {
			memberRoles[member.User.ID] = member.Roles
		}
		lastUserID = gm[len(gm)-1].User.ID
	}

	for userID, roles := range memberRoles {
		if _, exists := userSet[userID]; !exists {
			if _, err := g.userService.Join(usecase.UserJoinRequest{UserID: userID}); err != nil {
				g.logger.Warn("Failed to add user %s: %v", userID, err)
				continue
			}
		}

		for _, roleID := range roles {
			campus, ok := roleToCampus[roleID]
			if !ok {
				continue
			}

			if _, err := g.userService.SetUserCampus(usecase.SetUserCampusRequest{
				UserID: userID,
				Campus: campus,
			}); err != nil {
				g.logger.Warn("Failed to set user campus %s to %s: %v", userID, campus, err)
			}
		}

		delete(userSet, userID)
	}

	for userID := range userSet {
		if _, err := g.userService.Leave(usecase.UserLeaveRequest{UserID: userID}); err != nil {
			g.logger.Warn("Failed to remove user %s: %v", userID, err)
		}
	}

	g.logger.Info("Started discord gateway for application %s", g.applicationID)
}

func (g *gateway) Stop() {
	g.logger.Info("Stopped discord gateway")
}

func (g *gateway) RegisterCommand(def *discordgo.ApplicationCommand, handleFunc interaction.HandleFunc) {
	if _, ok := g.handleFuncs[def.Name]; ok {
		g.logger.Warn("Command %s already registered", def.Name)
		return
	}

	_, err := g.session.ApplicationCommandCreate(g.applicationID, "", def)
	if err != nil {
		g.logger.Error("Failed to register command %s: %v", def.Name, err)
		return
	}

	g.handleFuncs[def.Name] = handleFunc
	g.logger.Info("Registered command %s", def.Name)
}

func (g *gateway) RegisterComponent(def discordgo.MessageComponent, handleFunc interaction.HandleFunc) {
	var customID string
	switch v := def.(type) {
	case discordgo.Button:
		customID = v.CustomID
	case discordgo.SelectMenu:
		customID = v.CustomID
	}

	routingKey, _, _ := strings.Cut(customID, "?")
	if _, ok := g.handleFuncs[routingKey]; ok {
		g.logger.Warn("Component %s already registered", routingKey)
		return
	}

	g.handleFuncs[customID] = handleFunc
	g.logger.Info("Registered component %s", routingKey)
}

func (g *gateway) RegisterModal(def *discordgo.InteractionResponseData, handleFunc interaction.HandleFunc) {
	if _, ok := g.handleFuncs[def.CustomID]; ok {
		g.logger.Warn("Modal %s already registered", def.CustomID)
		return
	}

	g.handleFuncs[def.CustomID] = handleFunc
	g.logger.Info("Registered modal %s", def.CustomID)
}
