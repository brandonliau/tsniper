package main

import (
	"sync"
	"sync/atomic"
	"time"

	"tsniper/internal/config"
	"tsniper/pkg/logger"

	"github.com/bwmarrin/discordgo"
)

func purge() {
	logger := logger.NewStdLogger(logger.LevelInfo)

	cfg, err := config.NewYamlConfig("./config/config.yml")
	if err != nil {
		logger.Fatal("Failed to create config: %v", err)
	}

	s, err := discordgo.New("Bot " + cfg.Discord.Token)
	if err != nil {
		logger.Fatal("Failed to create discord session: %v", err)
	}
	s.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)

	s.AddHandler(func(_ *discordgo.Session, r *discordgo.RateLimit) {
		logger.Warn("Rate limited on %s, retry after %v", r.URL, r.RetryAfter)
	})

	err = s.Open()
	if err != nil {
		logger.Fatal("Failed to establish websocket connection: %v", err)
	}
	defer s.Close()

	botID := s.State.User.ID

	var members []*discordgo.Member
	afterID := ""
	for {
		batch, err := s.GuildMembers(cfg.Discord.GuildID, afterID, 1000)
		if err != nil {
			logger.Fatal("Failed to fetch guild members: %v", err)
		}
		if len(batch) == 0 {
			break
		}
		members = append(members, batch...)
		afterID = batch[len(batch)-1].User.ID
		if len(batch) < 1000 {
			break
		}
	}

	logger.Info("Found %d guild members", len(members))

	var totalDeleted atomic.Int64
	var wg sync.WaitGroup
	sem := make(chan struct{}, 5)

	for _, member := range members {
		if member.User.Bot {
			continue
		}

		wg.Add(1)
		sem <- struct{}{}
		go func(m *discordgo.Member) {
			defer wg.Done()
			defer func() { <-sem }()

			ch, err := s.UserChannelCreate(m.User.ID)
			if err != nil {
				logger.Warn("Failed to open DM channel for user %s: %v", m.User.ID, err)
				return
			}

			deleted := 0
			var beforeID string
			for {
				messages, err := s.ChannelMessages(ch.ID, 100, beforeID, "", "")
				if err != nil {
					logger.Error("Failed to fetch messages for user %s: %v", m.User.ID, err)
					break
				}
				if len(messages) == 0 {
					break
				}

				for _, msg := range messages {
					if msg.Author.ID == botID {
						err := s.ChannelMessageDelete(ch.ID, msg.ID)
						if err != nil {
							logger.Error("Failed to delete message %s for user %s: %v", msg.ID, m.User.ID, err)
						} else {
							deleted++
						}
						time.Sleep(750 * time.Millisecond)
					}
				}

				beforeID = messages[len(messages)-1].ID
				if len(messages) < 100 {
					break
				}
			}

			if deleted > 0 {
				logger.Info("Deleted %d messages for user %s", deleted, m.User.ID)
				totalDeleted.Add(int64(deleted))
			}
		}(member)
	}

	wg.Wait()
	logger.Info("Purged %d DM messages total", totalDeleted.Load())
}
