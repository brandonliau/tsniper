package component

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type RegisterButton struct {
	link string
}

func NewRegisterButton(index, term, year, campus string) *RegisterButton {
	url := fmt.Sprintf(
		"https://sims.rutgers.edu/webreg/editSchedule.htm?login=cas&semesterSelection=%s&indexList=%s&campus=%s",
		term+year,
		index,
		campus,
	)
	return &RegisterButton{
		link: url,
	}
}

func (c *RegisterButton) Component() discordgo.MessageComponent {
	return discordgo.Button{
		Label: "Register",
		Style: discordgo.LinkButton,
		URL:   c.link,
	}
}
