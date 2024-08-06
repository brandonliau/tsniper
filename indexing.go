package main

import (
	"fmt"
	"time"
)

var dayMap = map[string]string {
	"M": "Mon",
	"T": "Tue",
	"W": "Wed",	
	"H": "Thu",
	"F": "Fri",
	"S": "Sat",
	"U": "Sun",
}

var campusMap = map[string]string{
	"LIVINGSTON": ":orange_square:",
	"BUSCH": ":blue_square:",
	"COLLEGE AVENUE": ":yellow_square: ",
	"DOUGLAS/COOK": ":green_square:",
	"DOWNTOWN NB": "red_square",
	"OFF CAMPUS": ":white_large_square:",
}

func (runTimeData *RunTimeData) IndexCourses(campus string, season string) {
	var rawCourseData []RawCourseData = CoursesAPI(campus, runTimeData.Config.Seasons[season])
	runTimeData.AllCourses[campus + season] = make(map[string]struct{})
	lastOpens := runTimeData.GetLastOpens(campus, season)
	tx, _ := runTimeData.Db.Begin()
	for _, course := range rawCourseData {
		title := course.Title
		courseString := course.CourseString
		for _, section := range course.Sections {
			number := section.Section
			index := section.Index
			instructors := section.Instructors
			notes := section.SectionNotes
			meetingData := ""
			lastOpen := lastOpens[section.Index]
			if len(notes) == 0 {
				notes = "No notes for this section"
			}
			if section.Instructors == "" {
				instructors = "Unknown"
			}
			if section.OpenStatus {
				lastOpen = 0
			}
			for _, meeting := range section.MeetingTimes {
				if meeting.CampusLocation == "O" && meeting.PmCode == "" { // Async class with no meeting time
					meetingData += "Asynchronous Content :computer: Online\n"
				} else if meeting.CampusLocation == "O" && meeting.PmCode != "" { // Async class with set meeting time
					day := dayMap[meeting.MeetingDay]
					startTime := ParseTime(meeting.StartTimeMilitary)
					endTime := ParseTime(meeting.EndTimeMilitary)
					meetingData += fmt.Sprintf("%s %s - %s :computer: Online\n", day, startTime, endTime)
				} else if meeting.CampusLocation == "3" && meeting.PmCode == "" { // In person class with hours by arrangement
					campus := meeting.CampusName
					meetingData += fmt.Sprintf("Hours by arrangement %s %s\n", campusMap[campus], campus)
				} else if meeting.CampusLocation != "" && meeting.PmCode != "" { // In person class with set meeting time
					day := dayMap[meeting.MeetingDay]
					startTime := ParseTime(meeting.StartTimeMilitary)
					endTime := ParseTime(meeting.EndTimeMilitary)
					campus := meeting.CampusName
					var location string
					if campus == "** INVALID **" { // Unknown meeting location
						location = "Unknown"
					} else {
						location = meeting.BuildingCode + "-" + meeting.RoomNumber
					}
					meetingData += fmt.Sprintf("%s %s - %s %s **%s** (`%s`)\n",
						day, startTime, endTime, campusMap[campus], campus, location)
				} else {
					meetingData += "Hours by arrangement\n"
				}
			}
			runTimeData.AllCourses[campus + season][index] = struct{}{}
			query := fmt.Sprintf(
				"INSERT OR REPLACE INTO %s (course_index, title, course_string, section, instructors, notes, meeting, season, last_open)" +
				"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", campus)
			tx.Exec(query, index, title, courseString, number, instructors, notes, meetingData[:len(meetingData)-1], season, lastOpen)
		}
	}
	_ = tx.Commit()
}

func (runTimeData RunTimeData) IndexSync() {
	for _, campus := range runTimeData.Config.CurrentCampuses {
		for _, season := range runTimeData.Config.CurrentSeasons {
			runTimeData.IndexCourses(campus, season)
		}
	}
	fmt.Printf("SUCCESS @ %s : INDEX COURSES\n", time.Now().Format("2006-01-02 15:04:05.00000"))
}
