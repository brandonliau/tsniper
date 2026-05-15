package discord

import (
	"fmt"
	"time"

	"tsniper/internal/application/event"
	"tsniper/internal/application/view"
	"tsniper/internal/config"
	"tsniper/internal/interfaces/discord/component"
	"tsniper/internal/interfaces/discord/presentation"

	"tsniper/pkg/eventbus"
	"tsniper/pkg/logger"
	"tsniper/pkg/utils"

	"github.com/bwmarrin/discordgo"
)

type openCourseNotifier struct {
	stop            chan struct{}
	session         *discordgo.Session
	eventSubscriber eventbus.Subscriber[event.CourseOpen]
	customization   *config.CustomizationConfig
	logger          logger.Logger
}

func NewOpenCourseNotifier(eventSubscriber eventbus.Subscriber[event.CourseOpen], session *discordgo.Session, customization *config.CustomizationConfig, logger logger.Logger) *openCourseNotifier {
	return &openCourseNotifier{
		stop:            make(chan struct{}),
		session:         session,
		eventSubscriber: eventSubscriber,
		customization:   customization,
		logger:          logger,
	}
}

func (n *openCourseNotifier) Start() {
	ch := n.eventSubscriber.Subscribe()
	go func() {
		defer close(n.stop)
		for e := range ch {
			if err := n.execute(e); err != nil {
				n.logger.Error("Failed to send open course notification: %v", err)
			}
		}
	}()

	n.logger.Info("Started open course discord notifier")
}

func (n *openCourseNotifier) Stop() {
	<-n.stop
	n.logger.Info("Stopped open course discord notifier")
}

func (n *openCourseNotifier) execute(e event.CourseOpen) error {
	var embed *discordgo.MessageEmbed
	switch e.Type {
	case event.CourseOpenNotification:
		embed = n.notificationEmbed(e.Course)
	case event.CourseOpenAutosnipeSuccess:
		embed = n.successfulAutosnipeEmbed(e.Course)
	case event.CourseOpenAutosnipeFailed:
		embed = n.failedAutosnipeEmbed(e.Course, e.ErrMessage)
	}

	components := []discordgo.MessageComponent{
		component.RegisterComponentDefinition(e.Course.Index, e.Course.Campus, e.Course.Term, e.Course.Year),
		component.ResnipeComponentDefinition(
			utils.KeyValue[string, string]{Key: "index", Value: e.Course.Index},
			utils.KeyValue[string, string]{Key: "campus", Value: e.Course.Campus},
			utils.KeyValue[string, string]{Key: "term", Value: e.Course.Term},
			utils.KeyValue[string, string]{Key: "year", Value: e.Course.Year},
		),
	}

	for _, userID := range e.UserIDs {
		ch, err := n.session.UserChannelCreate(userID)
		if err != nil {
			n.logger.Error("Failed to create user channel: %v", err)
			continue
		}

		msg := &discordgo.MessageSend{
			Content: fmt.Sprintf("<@%s> ", userID),
			Embeds:  []*discordgo.MessageEmbed{embed},
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{Components: components},
			},
		}

		_, err = n.session.ChannelMessageSendComplex(ch.ID, msg)
		if err != nil {
			n.logger.Error("Failed to send message to user %s: %v", userID, err)
			continue
		}

		n.logger.Info("Notified user %s that %s is open", userID, e.Course.Index)
	}

	return nil
}

func (n *openCourseNotifier) notificationEmbed(crs *view.CourseView) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title: fmt.Sprintf("%s (Section %s) has opened!", crs.Title, crs.Section),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Course Name",
				Value:  fmt.Sprintf("`%s`", crs.Title),
				Inline: true,
			},
			{
				Name:   "Index",
				Value:  fmt.Sprintf("`%s`", crs.Index),
				Inline: true,
			},
			{
				Name:   "Section",
				Value:  fmt.Sprintf("`%s`", crs.Section),
				Inline: true,
			},
		},
		Color: presentation.Blue,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: n.customization.Thumbnail,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: time.Now().Format("01/02/2006 03:04:05 PM"),
		},
	}
}

func (n *openCourseNotifier) successfulAutosnipeEmbed(crs *view.CourseView) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title: fmt.Sprintf("Autosniped %s!", crs.Title),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Course Name",
				Value:  fmt.Sprintf("`%s`", crs.Title),
				Inline: true,
			},
			{
				Name:   "Index",
				Value:  fmt.Sprintf("`%s`", crs.Index),
				Inline: true,
			},
			{
				Name:   "Section",
				Value:  fmt.Sprintf("`%s`", crs.Section),
				Inline: true,
			},
		},
		Color: presentation.Blue,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: n.customization.Thumbnail,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: time.Now().Format("01/02/2006 03:04:05 PM"),
		},
	}
}

func (n *openCourseNotifier) failedAutosnipeEmbed(crs *view.CourseView, message string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title: fmt.Sprintf("Failed to autosnipe %s!", crs.Title),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Course Name",
				Value:  fmt.Sprintf("`%s`", crs.Title),
				Inline: true,
			},
			{
				Name:   "Index",
				Value:  fmt.Sprintf("`%s`", crs.Index),
				Inline: true,
			},
			{
				Name:   "Section",
				Value:  fmt.Sprintf("`%s`", crs.Section),
				Inline: true,
			},
			{
				Name:   "Reason",
				Value:  fmt.Sprintf("`%s`", message),
				Inline: false,
			},
		},
		Color: presentation.Red,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: n.customization.Thumbnail,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: time.Now().Format("01/02/2006 03:04:05 PM"),
		},
	}
}
