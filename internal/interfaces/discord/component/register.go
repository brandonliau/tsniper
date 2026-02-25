package component

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func RegisterComponentDefinition(index string, campus string, term string, year string) discordgo.Button {
	url := fmt.Sprintf(
		"https://sims.rutgers.edu/webreg/editSchedule.htm?login=cas&semesterSelection=%s%s&indexList=%s",
		term,
		year,
		index,
	)

	return discordgo.Button{
		Label: "Register",
		Style: discordgo.LinkButton,
		URL:   url,
	}
}
