package manager

import (
	"Tsniper/internal/command"
	"Tsniper/internal/component"
	"Tsniper/internal/notifier"
	"Tsniper/internal/repository"
	"Tsniper/internal/shared"
	"Tsniper/pkg/config"
	"Tsniper/pkg/logger"

	"github.com/bwmarrin/discordgo"
)

type mockManager struct {
	dCfg            *config.DiscordConfig
	session           *discordgo.Session
	application       *discordgo.Application
	retrievedCommands []*discordgo.ApplicationCommand
	repo              repository.Repository
	notifier          notifier.Notifier
	logger            logger.Logger
	commands          map[string]command.Command
	components        map[string]component.Component
}

func NewMockManager(dCfg *config.DiscordConfig, s *discordgo.Session, repo repository.Repository, logger logger.Logger, notifier notifier.Notifier) *mockManager {
	return &mockManager{
		dCfg:            dCfg,
		session:           s,
		retrievedCommands: make([]*discordgo.ApplicationCommand, 0),
		repo:              repo,
		notifier:          notifier,
		logger:            logger,
		commands:          make(map[string]command.Command),
		components:        make(map[string]component.Component),
	}
}

func (m *mockManager) RetreiveCommands() {
	var err error
	m.retrievedCommands, err = m.session.ApplicationCommands(m.application.ID, "")
	if err != nil {
		m.logger.Error("Failed to retrieve application commands")
	}
}

func (m *mockManager) RegisterCommand(c command.Command) {
	cname := c.Command().Name
	if _, ok := m.commands[cname]; ok {
		m.logger.Warn("Application command %s already registered", cname)
	}
	for _, ccmd := range m.retrievedCommands {
		if ccmd.Name == cname {
			m.repo.RegisterCommand(ccmd.Name, ccmd.ID)
			m.commands[cname] = c
			m.logger.Debug("Registered command %v", cname)
			return
		}
	}
	ccmd, err := m.session.ApplicationCommandCreate(m.session.State.User.ID, "", c.Command())
	if err != nil {
		m.logger.Error("Failed to add application command %s : %v", cname, err)
	}
	m.repo.RegisterCommand(ccmd.Name, ccmd.ID)
	m.commands[cname] = c
	m.logger.Debug("Registered command %v", cname)
}

func (m *mockManager) RegisterComponent(c component.Component) {
	cname := c.CustomID()
	if _, ok := m.components[cname]; ok {
		m.logger.Warn("Application component %s already registered", cname)
	}
	m.repo.RegisterComponent(cname, c.Component())
	m.components[cname] = c
	m.logger.Debug("Registered component %v", cname)
}

func (m *mockManager) CommandInteractionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var userID string
	if i.Member != nil {
		userID = i.Member.User.ID
	} else {
		userID = i.User.ID
	}
	cmdArgs := &shared.CmdArgs{
		Session:     s,
		Interaction: i,
		UserID:      userID,
	}

	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	var ir *discordgo.InteractionResponse
	var err error
	command := m.commands[i.ApplicationCommandData().Name]
	ir, err = command.Execute(cmdArgs)
	if err != nil {
		m.logger.Error("Failed to execute /%s: %v", command.Command().Name, err)
	}
	m.logger.Debug("%s executed /%s", userID, i.ApplicationCommandData().Name)

	err = m.notifier.SendResponse(i, ir)
	if err != nil {
		m.logger.Error("Failed to respond to user %s: %v", userID, err)
	}
}

func (m *mockManager) ComponentInteractionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var userID string
	if i.Member != nil {
		userID = i.Member.User.ID
	} else {
		userID = i.User.ID
	}
	cmdArgs := &shared.CmdArgs{
		Session:     s,
		Interaction: i,
		UserID:      userID,
	}

	if i.Type != discordgo.InteractionMessageComponent {
		return
	}

	var ir *discordgo.InteractionResponse
	var err error
	component := m.components[i.MessageComponentData().CustomID]
	ir, err = component.Execute(cmdArgs)
	if err != nil {
		m.logger.Error("Failed to execute /%s: %v", component.CustomID(), err)
		return
	}
	m.logger.Debug("%s executed /%s", userID, component.CustomID())

	err = m.notifier.SendResponse(i, ir)
	if err != nil {
		m.logger.Error("Failed to respond to user %s: %v", userID, err)
	}
}
