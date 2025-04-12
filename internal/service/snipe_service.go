package service

import (
	"fmt"
	"net/http"
	"time"

	"tsniper/internal/component"
	"tsniper/internal/notifier"
	"tsniper/internal/repository"
	"tsniper/internal/shared"
	"tsniper/pkg/config"
	"tsniper/pkg/database"
	"tsniper/pkg/logger"
	"tsniper/pkg/multitick"

	"github.com/bwmarrin/discordgo"
)

type snipeService struct {
	sCfg     *config.ServiceConfig
	dCfg     *config.DiscordConfig
	session  *discordgo.Session
	client   *http.Client
	ticker   *multitick.Ticker
	stop     chan bool
	repo     repository.Repository
	notifier notifier.Notifier
	db       database.Database
	logger   logger.Logger
}

func NewSnipeService(
	sCfg *config.ServiceConfig,
	dCfg *config.DiscordConfig,
	session *discordgo.Session,
	repo repository.Repository,
	nofitifer notifier.Notifier,
	db database.Database,
	logger logger.Logger,
) *snipeService {
	transport := &http.Transport{
		MaxIdleConns:        6,
		MaxIdleConnsPerHost: 6,
		IdleConnTimeout:     90 * time.Second,
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}
	snipeService := &snipeService{
		sCfg:     sCfg,
		dCfg:     dCfg,
		session:  session,
		client:   client,
		stop:     make(chan bool),
		repo:     repo,
		notifier: nofitifer,
		db:       db,
		logger:   logger,
	}
	err := snipeService.migrate()
	if err != nil {
		logger.Fatal("Failed to migrate snipe tables: %v", err)
	}
	return snipeService
}

func (s *snipeService) migrate() error {
	err := s.db.ExecSQLFile("./pkg/database/migrations/snipe.sql")
	return err
}

func (s *snipeService) Start() error {
	s.logger.Info("Started snipe service")
	err := s.migrate()
	if err != nil {
		return err
	}
	s.ticker = multitick.NewMultiTicker(time.Second, 500*time.Millisecond)
	go s.snipeLoop()
	return nil
}

func (s *snipeService) Stop() error {
	s.logger.Info("Stopped snipe service")
	s.ticker.Stop()
	s.stop <- true
	return nil
}

func (s *snipeService) snipeCheck(openChan chan []string, season, year, term, campus string) {
	c := s.ticker.Subscribe()
	prevOpened := OpenSections(s.client, year, term, campus)
	for {
		select {
		case <-c:
			openSections := OpenSections(s.client, year, term, campus)
			trackedIndices := s.repo.TrackedIndices(campus, season)
			openCourses := shared.Intersection(openSections, trackedIndices)
			if len(openCourses) > 0 {
				openCourses = append(openCourses, campus+season) // add [campus + season] to the end for processing loop
				openChan <- openCourses                          // send the data to the processing loop for it to notify users
			}
			openDiff := shared.Difference(openSections, prevOpened)
			closeDiff := shared.Difference(prevOpened, openSections)

			for _, index := range openDiff {
				s.repo.UpdateLastOpen(index, campus, season, 0)
			}
			for _, index := range closeDiff {
				s.repo.UpdateLastOpen(index, campus, season, time.Now().Unix())
			}
			prevOpened = openSections
		case <-s.stop:
			return
		}
	}
}

func (s *snipeService) snipeLoop() {
	openChan := make(chan []string)
	for _, campus := range s.sCfg.Campuses {
		for _, season := range s.sCfg.Seasons {
			year := s.sCfg.SeasonData[season].Year
			term := s.sCfg.SeasonData[season].Term
			go s.snipeCheck(openChan, season, year, term, campus)
		}
	}
	go func() {
		<-s.stop
		close(openChan)
	}()
	for data := range openChan {
		campusSeason := data[len(data)-1]
		campus := campusSeason[0:2]
		season := campusSeason[2:]
		sections := data[:len(data)-1]
		for _, index := range sections {
			course := s.repo.CourseEntry(index, campus, season)
			embed := s.snipeEmbed(course)
			year := s.sCfg.SeasonData[season].Year
			term := s.sCfg.SeasonData[season].Term
			registerButton := component.NewRegisterButton(index, term, year, campus).Component()
			resnipeButton := s.repo.RetrieveComponents("resnipe")[0]
			for _, userID := range s.repo.SnipeUsers(index, campus, season) {
				user, _ := s.session.User(userID)
				dmChannel, err := s.notifier.CreateDMChannel(userID)
				if err != nil {
					s.logger.Warn("Failed to create dm channel with user %s", userID)
					continue
				}
				err = s.notifier.SendChannelMessage(dmChannel, notifier.MessageSend(user.Mention(), embed, registerButton, resnipeButton))
				if err != nil {
					s.logger.Warn("Failed to message user %s: %v", userID, err)
					s.repo.RemoveSnipe(userID, index, campus, season)
					continue
				} else {
					s.logger.Info("Notified user %s that %s is open", userID, index)
				}
				s.repo.RemoveSnipe(userID, index, campus, season)
				s.repo.Remove(index, campus, season)
			}
		}
	}
}

func (s *snipeService) snipeEmbed(course shared.CourseEntry) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title: fmt.Sprintf("%s (Section %s) has opened!", course.Title, course.Section),
		Color: shared.Blue,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Course Name",
				Value:  fmt.Sprintf("`%s`", course.Title),
				Inline: true,
			},
			{
				Name:   "Index",
				Value:  fmt.Sprintf("`%s`", course.Index),
				Inline: true,
			},
			{
				Name:   "Section",
				Value:  fmt.Sprintf("`%s`", course.Section),
				Inline: true,
			},
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: s.dCfg.Image,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: time.Now().Format("01/02/2006 03:04:05 PM"),
		},
	}
}
