package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	
	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron/v3"
	_ "modernc.org/sqlite"
	
	"encoding/json"
)

var s *discordgo.Session
var snipeTicker = time.NewTicker(175 * time.Millisecond)
var snipeClose = make(chan bool)
var indexingTimer = cron.New()

func main() {
	var runTimeData RunTimeData
	runTimeData.LoadConfig("./config.yml")
	runTimeData.InitRunTimeData()
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
	// registered, _ := s.ApplicationCommandBulkOverwrite(s.State.User.ID, "", Commands)
	fmt.Printf("SUCCESS @ %s : REGISTER ALL COMMANDS\n", time.Now().Format("2006-01-02 15:04:05.00000"))

	// SaveCommands(registered, "commands.json")
	registered := LoadCommands("./commands.json")

	runTimeData.Registered = make(map[string]string)
	for _, command := range registered {
		temp := *command
		runTimeData.Registered[temp.Name] = temp.ID
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	snipeTicker.Stop()
	indexingTimer.Stop()
	snipeClose <- true
	runTimeData.Db.Close()

	_ = s.Close()
	fmt.Printf("SUCCESS @ %s : CLOSE WEBSOCKET CONNECTION\n", time.Now().Format("2006-01-02 15:04:05.00000"))
	// _, _ = s.ApplicationCommandBulkOverwrite(s.State.User.ID, "", nil)
	// fmt.Printf("SUCCESS @ %s : REMOVE ALL COMMANDS\n", time.Now().Format("2006-01-02 15:04:05.00000"))
}

func SaveCommands(commands []*discordgo.ApplicationCommand, filename string) {
	updatedJson, _ := json.MarshalIndent(commands, "", "  ")
	os.WriteFile(filename, updatedJson, 0644)
}

func LoadCommands(filename string) []*discordgo.ApplicationCommand {
	var registered = []*discordgo.ApplicationCommand{}
	rawJson, _ := os.ReadFile(filename)
	_ = json.Unmarshal(rawJson, &registered)
	return registered
}
