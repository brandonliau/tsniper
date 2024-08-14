package debug

import (
	"encoding/json"
	"os"
	
	"github.com/bwmarrin/discordgo"
)

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
