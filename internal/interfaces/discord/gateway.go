package discord

import (
	"strings"
	"tsniper/internal/application/usecase"
	"tsniper/internal/config"
	"tsniper/internal/interfaces/discord/interaction"

	"tsniper/pkg/logger"
	"tsniper/pkg/utils"

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
	var users []string
	for _, usr := range res.Users {
		users = append(users, usr.ID)
	}

	var members []string
	var lastUserID string
	for {
		gm, _ := g.session.GuildMembers(g.guildID, lastUserID, 1000)
		if len(gm) < 1 {
			break
		}
		for _, member := range gm {
			members = append(members, member.User.ID)
		}
		lastUserID = gm[len(gm)-1].User.ID
	}

	toAdd := utils.Difference(members, users)
	toRemove := utils.Difference(users, members)

	for _, userID := range toAdd {
		_, err := g.userService.Join(usecase.UserJoinRequest{UserID: userID})
		if err != nil {
			g.logger.Warn("Failed to add user %s: %v", userID, err)
			continue
		}

		member, err := g.session.GuildMember(g.guildID, userID)
		if err != nil {
			g.logger.Warn("Failed to get member %s: %v", userID, err)
			continue
		}

		for _, roleID := range member.Roles {
			roleName, ok := g.cfg.Roles[roleID]
			if !ok {
				continue
			}

			req := usecase.SetUserCampusRequest{
				UserID: userID,
				Campus: roleName,
			}

			_, err := g.userService.SetUserCampus(req)
			if err != nil {
				g.logger.Warn("Failed to set user campus %s to %s: %v", userID, roleName, err)
				continue
			}
		}
	}

	for _, userID := range toRemove {
		_, err := g.userService.Leave(usecase.UserLeaveRequest{UserID: userID})
		if err != nil {
			g.logger.Warn("Failed to remove user %s: %v", userID, err)
			continue
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
