package manager

import (
	"fmt"
	"time"

	"Tsniper/internal/notifier"
	"Tsniper/internal/shared"

	"github.com/bwmarrin/discordgo"
)

func (m *mockManager) ReadyHandler(s *discordgo.Session, r *discordgo.Ready) {
	s.UpdateCustomStatus("👁️‍🗨️ Monitoring...")
	m.application = r.Application
	m.logger.Info("Ready event")
}

func (m *mockManager) ResumedHandler(s *discordgo.Session, r *discordgo.Resumed) {
	m.logger.Info("Resumed event")
}

func (m *mockManager) RateLimitHandler(s *discordgo.Session, r *discordgo.RateLimit) {
	m.logger.Info("Rate limit event")
}

func (m *mockManager) GuildMemberAddHandler(s *discordgo.Session, g *discordgo.GuildMemberAdd) {
	embed := m.joinEmbed(g.Member.User)
	m.notifier.SendChannelMessage(m.dCfg.Boarding, notifier.MessageSend("", embed))
	m.logger.Info("Guild member add event")
}

func (m *mockManager) GuildMemberRemoveHandler(s *discordgo.Session, g *discordgo.GuildMemberRemove) {
	snipes := m.repo.Snipes(g.Member.User.ID)
	m.repo.ClearSnipe(g.Member.User.ID)
	for _, snipe := range snipes {
		index, campus, season := snipe[0], snipe[1], snipe[2]
		m.repo.Remove(index, campus, season)
	}
	m.logger.Info("Guild member remove event")
}

func (m *mockManager) joinEmbed(user *discordgo.User) *discordgo.MessageEmbed {
	guild, _ := m.session.State.Guild(m.dCfg.Guild)
	return &discordgo.MessageEmbed{
		Title:       "Welcome to the TSniper server!",
		Description: fmt.Sprintf("%s has joined the server!\n\nYou are user **#%d**!", user.Mention(), guild.MemberCount),
		Color:       shared.Green,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: user.AvatarURL(""),
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: time.Now().Format("01/02/2006 03:04:05 PM"),
		},
	}
}
