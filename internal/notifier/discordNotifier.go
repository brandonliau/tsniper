package notifier

import (
	"github.com/bwmarrin/discordgo"
)

type discordNotifier struct {
	session *discordgo.Session
}

func NewDiscordNotifier(s *discordgo.Session) *discordNotifier {
	return &discordNotifier{
		session: s,
	}
}

func (n *discordNotifier) SendResponse(i *discordgo.InteractionCreate, ir *discordgo.InteractionResponse) error {
	err := n.session.InteractionRespond(i.Interaction, ir)
	if err != nil {
		return err
	}
	return nil
}

func (n *discordNotifier) SendComplexMessage(channelID string, data *discordgo.MessageSend) error {
	_, err := n.session.ChannelMessageSendComplex(
		channelID,
		data,
	)
	if err != nil {
		return err
	}
	return nil
}

func (n *discordNotifier) CreateDMChannel(userID string) (string, error) {
	dmChannel, err := n.session.UserChannelCreate(userID)
	if err != nil {
		return "", err
	}
	return dmChannel.ID, nil
}
