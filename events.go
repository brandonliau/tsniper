package main

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

func (runTimeData *RunTimeData) OnReady(s *discordgo.Session, r *discordgo.Ready) {
	runTimeData.IndexSync()
	runTimeData.SyncUsers()
	runTimeData.SyncTracking()
	fmt.Printf("SUCCESS @ %s : SYNC TRACKING\n", time.Now().Format("2006-01-02 15:04:05.00000"))
	runTimeData.StartTime = time.Now().Unix()
	indexingTimer.AddFunc("0 1 * * *", runTimeData.IndexSync)
	indexingTimer.AddFunc("0 13 * * *", runTimeData.IndexSync)
	go runTimeData.SnipeCheck()
	indexingTimer.Start()
	s.UpdateCustomStatus("👁️‍🗨️ Monitoring...")
	fmt.Println("***************************** BOT RUNNING *****************************")
}

func (runTimeData *RunTimeData) InteractionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		if command, ok := CommandHandlers[i.ApplicationCommandData().Name]; ok {
			command(s, i)
		}
	case discordgo.InteractionMessageComponent:
		if command, ok := ComponentHandlers[i.MessageComponentData().CustomID]; ok {
			command(s, i)
		}
	}
}

func (runTimeData *RunTimeData) OnMemberJoin(s *discordgo.Session, g *discordgo.GuildMemberAdd) {
	embed := runTimeData.JoinEmbed(g.Member.User)
	s.ChannelMessageSendComplex(
		runTimeData.Config.Boarding,
		&discordgo.MessageSend{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	)
}

func (runTimeData *RunTimeData) OnMemberRemove(s *discordgo.Session, g *discordgo.GuildMemberRemove) {
	runTimeData.ClearSnipe(g.Member.User.ID)
	runTimeData.SyncTracking()
}
