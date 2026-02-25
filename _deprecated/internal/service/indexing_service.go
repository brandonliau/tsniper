package service

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"tsniper/internal/repository"
	"tsniper/pkg/config"
	"tsniper/pkg/database"
	"tsniper/pkg/logger"

	"github.com/robfig/cron/v3"
)

var dayMap = map[string]string{
	"M": "Mon",
	"T": "Tue",
	"W": "Wed",
	"H": "Thu",
	"F": "Fri",
	"S": "Sat",
	"U": "Sun",
}

var campusMap = map[string]string{
	"LIVINGSTON":     ":orange_square:",
	"BUSCH":          ":blue_square:",
	"COLLEGE AVENUE": ":yellow_square: ",
	"DOUGLAS/COOK":   ":green_square:",
	"DOWNTOWN NB":    "red_square",
	"OFF CAMPUS":     ":white_large_square:",
}

type indexingService struct {
	config *config.ServiceConfig
	client *http.Client
	cron   *cron.Cron
	repo   repository.Repository
	db     database.Database
	logger logger.Logger
}

func NewIndexingService(cfg *config.ServiceConfig, repo repository.Repository, db database.Database, logger logger.Logger) *indexingService {
	transport := &http.Transport{
		MaxIdleConns:        6,
		MaxIdleConnsPerHost: 6,
		IdleConnTimeout:     90 * time.Second,
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   15 * time.Second,
	}
	indexingService := &indexingService{
		config: cfg,
		client: client,
		cron:   cron.New(),
		repo:   repo,
		db:     db,
		logger: logger,
	}

	err := indexingService.migrate()
	if err != nil {
		logger.Fatal("Failed to migrate indexing tables: %v", err)
	}
	return indexingService
}

func (s *indexingService) migrate() error {
	err := s.db.ExecSQLFile("./pkg/database/migrations/indexing.sql")
	return err
}

func (s *indexingService) Start() error {
	s.logger.Info("Started indexing service")
	s.indexCoursesLoop()
	_, err := s.cron.AddFunc("0 1 * * *", s.indexCoursesLoop)
	if err != nil {
		return err
	}
	_, err = s.cron.AddFunc("0 13 * * *", s.indexCoursesLoop)
	if err != nil {
		return err
	}
	s.cron.Start()
	return nil
}

func (s *indexingService) Stop() error {
	s.logger.Info("Stopped indexing service")
	s.cron.Stop()
	return nil
}

func parseTime(time string) string {
	rawHour, _ := strconv.Atoi(time[:2])
	var hour int = 12
	if rawHour%12 != 0 {
		hour = rawHour % 12
	}
	timeString := strconv.Itoa(hour) + ":" + time[2:]
	if rawHour/12 < 1 {
		timeString += "AM"
	} else {
		timeString += "PM"
	}
	return timeString
}

func (s *indexingService) indexCourses(campus, season, year, term string) {
	start := time.Now()
	courses, err := Courses(s.client, year, term, campus)
	if err != nil {
		s.logger.Fatal("Failed to index courses: %v", err)
	}

	query := `INSERT OR REPLACE INTO courses 
		(course_index, title, course_string, section, instructors, notes, meeting, campus, season, last_open) 
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	stmt, err := s.db.PrepareExec(query)
	if err != nil {
		s.logger.Fatal("Failed to prepare indexing query: %v", err)
	}

	lastOpens := s.repo.LastOpens(campus, season)
	for _, course := range courses {
		title := course.Title
		courseString := course.CourseString
		for _, section := range course.Sections {
			var builder strings.Builder
			instructors := section.Instructors
			notes := section.SectionNotes
			lastOpen := lastOpens[section.Index]
			if section.Instructors == "" {
				instructors = "Unknown"
			}
			if len(notes) == 0 {
				notes = "No notes for this section"
			}
			if section.OpenStatus {
				lastOpen = 0
			} else if !section.OpenStatus && lastOpen == 0 {
				lastOpen = -1
			}
			for _, meeting := range section.MeetingTimes {
				if meeting.CampusLocation == "O" && meeting.PmCode == "" { // Async class with no meeting time
					builder.WriteString("Asynchronous Content :computer: Online\n")
				} else if meeting.CampusLocation == "O" && meeting.PmCode != "" { // Async class with set meeting time
					startTime := parseTime(meeting.StartTimeMilitary)
					endTime := parseTime(meeting.EndTimeMilitary)
					builder.WriteString(fmt.Sprintf("%s %s - %s :computer: Online\n", dayMap[meeting.MeetingDay], startTime, endTime))
				} else if meeting.CampusLocation == "3" && meeting.PmCode == "" { // In person class with hours by arrangement
					campus := meeting.CampusName
					builder.WriteString(fmt.Sprintf("Hours by arrangement %s %s\n", campusMap[campus], campus))
				} else if meeting.CampusLocation != "" && meeting.PmCode != "" { // In person class with set meeting time
					startTime := parseTime(meeting.StartTimeMilitary)
					endTime := parseTime(meeting.EndTimeMilitary)
					campus := meeting.CampusName
					var location string
					if campus == "** INVALID **" { // Unknown meeting location
						location = "Unknown"
					} else {
						location = meeting.BuildingCode + "-" + meeting.RoomNumber
					}
					builder.WriteString(fmt.Sprintf("%s %s - %s %s **%s** (`%s`)\n",
						dayMap[meeting.MeetingDay], startTime, endTime, campusMap[campus], campus, location))
				} else {
					builder.WriteString("Hours by arrangement\n")
				}
			}
			meetingData := builder.String()[:len(builder.String())-1]
			stmt.Exec(section.Index, title, courseString, section.Section, instructors, notes, meetingData, campus, season, lastOpen)
		}
	}
	s.logger.Info("Indexed %d courses in %v", len(courses), time.Since(start))
}

func (s *indexingService) indexCoursesLoop() {
	start := time.Now()
	s.db.Begin()
	for _, campus := range s.config.Campuses {
		for _, season := range s.config.Seasons {
			year := s.config.SeasonData[season].Year
			term := s.config.SeasonData[season].Term
			s.indexCourses(campus, season, year, term)
		}
	}
	s.db.Commit()
	s.logger.Info("Indexed %d campuses and seasons combinations in %v", len(s.config.Campuses)*len(s.config.Seasons), time.Since(start))
}
