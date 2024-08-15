package main

import (
	"database/sql"
)

type RawCourseData struct {
	Title              string `json:"title"`
	CourseString       string `json:"courseString"`
	Sections     []struct {
		Index				string `json:"index"`
		Section				string `json:"number"`
		Instructors 		string `json:"instructorsText"`
		SectionNotes        string `json:"sectionNotes"`
		OpenStatus			bool   `json:"openStatus"`
		MeetingTimes    []struct {
			CampusLocation    string `json:"campusLocation"`
			CampusName        string `json:"campusName"`
			PmCode            string `json:"pmCode"`
			MeetingDay        string `json:"meetingDay"`
			BuildingCode      string `json:"buildingCode"`
			RoomNumber        string `json:"roomNumber"`
			StartTimeMilitary string `json:"startTimeMilitary"`
			EndTimeMilitary   string `json:"endTimeMilitary"`
		} `json:"meetingTimes"`
	} `json:"sections"`
}

type CourseData struct {
	Title        string   `json:"title"`
	CourseString string   `json:"courseString"`
	Index		 string   `json:"index"`
	Section      string   `json:"section"`
	Instructors  string   `json:"instructors"`
	Notes        string   `json:"notes"`
	Meeting      string `json:"meeting"`
}

type Season struct {
	Name   string	`yaml:"name"`
	Term   string	`yaml:"term"`
	Year   string	`yaml:"year"`
}

type Config struct {
	Token             string			`yaml:"token"`
	Guild			  string 			`yaml:"guild"`
	Boarding		  string			`yaml:"boarding"`
	Image             string			`yaml:"image"`
	CurrentCampuses   []string			`yaml:"current_campuses"`
	CurrentSeasons    []string			`yaml:"current_seasons"`
	Seasons			  map[string]Season `yaml:"seasons"`
}

type RunTimeData struct {
	Config			 Config
	Db				 *sql.DB
	Tracking		 map[string]map[string]int
	Registered		 map[string]string
	StartTime		 int64
}
