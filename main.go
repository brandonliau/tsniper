package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron/v3"
	_ "modernc.org/sqlite"
)

var s *discordgo.Session
var snipeTicker = time.NewTicker(time.Second)
var snipeClose = make(chan bool)
var indexingTimer = cron.New()

func main() {
	var runTimeData RunTimeData
	runTimeData.InitRunTimeData("./tsniperconfig.json")
	runTimeData.InitCommands()
	runTimeData.InitHandlers()
	s, _ = discordgo.New("Bot " + runTimeData.Config.Token)
	s.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)
	s.AddHandler(runTimeData.OnReady)
	s.AddHandler(runTimeData.InteractionHandler)
	s.AddHandler(runTimeData.OnMemberJoin)
	s.AddHandler(runTimeData.OnMemberRemove)
	_ = s.Open()
	fmt.Printf("SUCCESS @ %s : ESTABLISH WEBSOCKET CONNECTION\n", time.Now().Format("2006-01-02 15:04:05.00000"))
	registered, _ := s.ApplicationCommandBulkOverwrite(s.State.User.ID, "", Commands)
	
	// SaveCommands(registered)
	// registered := LoadCommands()

	runTimeData.Registered = make(map[string]string)
	for _, command := range registered {
		temp := *command
		runTimeData.Registered[temp.Name] = temp.ID
	}
	fmt.Printf("SUCCESS @ %s : REGISTER ALL COMMANDS\n", time.Now().Format("2006-01-02 15:04:05.00000"))
	
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	
	snipeTicker.Stop()
	snipeClose <- true
	indexingTimer.Stop()
	runTimeData.Db.Close()

	_ = s.Close()
	fmt.Printf("SUCCESS @ %s : CLOSE WEBSOCKET CONNECTION\n", time.Now().Format("2006-01-02 15:04:05.00000"))
	_, _ = s.ApplicationCommandBulkOverwrite(s.State.User.ID, "", nil)
	fmt.Printf("SUCCESS @ %s : REMOVE ALL COMMANDS\n", time.Now().Format("2006-01-02 15:04:05.00000"))
}

func SaveCommands(commands []*discordgo.ApplicationCommand) {
	updatedJson, _ := json.MarshalIndent(commands, "", "  ")
	os.WriteFile("commands.json", updatedJson, 0644)
}

func LoadCommands() []*discordgo.ApplicationCommand {
	var registered = []*discordgo.ApplicationCommand{}
	rawJson, _ := os.ReadFile("./commands.json")
	_ = json.Unmarshal(rawJson, &registered)
	return registered
}