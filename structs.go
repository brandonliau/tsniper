package main

import (
	"database/sql"
	"sync"
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
	Name   string	`json:"name"`
	Term   string	`json:"term"`
	Year   string	`json:"year"`
	Emoji  string	`json:"emoji"`
}

type Config struct {
	Token             string			`json:"TOKEN"`
	Guild			  string 			`json:"GUILD"`
	Image             string			`json:"IMAGE"`
	Boarding		  string			`json:"BOARDING"`
	CurrentCampuses   []string			`json:"CURRENT_CAMPUSES"`
	CurrentSeasons    []string			`json:"CURRENT_SEASONS"`
	Seasons			  map[string]Season `json:"SEASONS"`
}

type RunTimeData struct {
	mtx				 *sync.Mutex
	Config			 Config
	Db				 *sql.DB
	Tracking		 map[string]map[string]int
	AllCourses		 map[string]map[string]struct{}
	PrevOpened		 map[string][]string
	Registered		 map[string]string
	StartTime		 int64
}
