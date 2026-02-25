package discord

import (
	"tsniper/internal/application/usecase"
	"tsniper/internal/interfaces/discord/interaction"

	"github.com/bwmarrin/discordgo"
)

func (g *gateway) InteractionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var customID string
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		customID = i.ApplicationCommandData().Name
	case discordgo.InteractionMessageComponent:
		customID = i.MessageComponentData().CustomID
	case discordgo.InteractionModalSubmit:
		customID = i.ModalSubmitData().CustomID
	default:
		return
	}

	routingKey, _ := interaction.DecodeCustomID(customID)
	handleFunc, ok := g.handleFuncs[routingKey]
	if !ok {
		g.logger.Error("%s interaction handler not found", routingKey)
		return
	}

	rsp, err := handleFunc(s, i)
	if err != nil {
		g.logger.Error("Failed to execute %s interaction handler: %v", routingKey, err)
		rsp = interaction.InteractionInitialResponse(
			interaction.WithContent("Something went wrong!"),
		)
	}

	if rsp != nil {
		err = g.session.InteractionRespond(i.Interaction, rsp)
		if err != nil {
			g.logger.Error("Failed to send interaction response: %v", err)
			return
		}
	}

	g.logger.Debug("Successfully executed %s interaction handler", routingKey)
}

func (g *gateway) ReadyHandler(s *discordgo.Session, r *discordgo.Ready) {
	g.session.UpdateCustomStatus("🎯 Locked & Loaded")
	g.logger.Info("Ready event")
}

func (g *gateway) ResumedHandler(s *discordgo.Session, r *discordgo.Resumed) {
	g.logger.Info("Resumed event")
}

func (g *gateway) RateLimitHandler(s *discordgo.Session, r *discordgo.RateLimit) {
	g.logger.Info("Rate limit event")
}

func (g *gateway) MemberJoinHandler(s *discordgo.Session, r *discordgo.GuildMemberAdd) {
	_, err := g.userService.Join(usecase.UserJoinRequest{UserID: r.User.ID})
	if err != nil {
		g.logger.Error("Failed to add user on join %s: %v", r.User.ID, err)
		return
	}
	g.logger.Info("Member join event")
}

func (g *gateway) MemberLeaveHandler(s *discordgo.Session, r *discordgo.GuildMemberRemove) {
	_, err := g.userService.Leave(usecase.UserLeaveRequest{UserID: r.User.ID})
	if err != nil {
		g.logger.Error("Failed to remove user on leave %s: %v", r.User.ID, err)
		return
	}
	g.logger.Info("Member leave event")
}
