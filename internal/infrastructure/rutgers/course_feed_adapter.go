package rutgers

import (
	"fmt"
	"strconv"
	"strings"

	"tsniper/internal/domain/course"
	"tsniper/internal/domain/scope"
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
	"COLLEGE AVENUE": ":yellow_square:",
	"DOUGLAS/COOK":   ":green_square:",
	"DOWNTOWN NB":    ":red_square:",
	"OFF CAMPUS":     ":white_large_square:",
}

func courseFeedDataToCourses(data []courseFeedData, scp scope.AcademicScope) []*course.Course {
	var courses []*course.Course
	for _, crs := range data {
		for _, sec := range crs.Sections {
			instructors := sec.Instructors
			if instructors == "" {
				instructors = "Unknown"
			}

			notes := sec.SectionNotes
			if notes == "" {
				notes = "No notes for this section"
			}

			var lines []string
			for _, m := range sec.MeetingTimes {
				lines = append(lines, courseFeedMeetingToMeeting(m))
			}
			meetings := strings.Join(lines, "\n")

			crs := course.NewCourse(sec.Index, crs.Title, crs.CourseString, sec.Section, instructors, notes, meetings, scp)
			
			if sec.OpenStatus {
				crs.Open()
			}

			courses = append(courses, crs)
		}
	}
	return courses
}

// todo: move this to discord search command handler once more interfaces are added
func courseFeedMeetingToMeeting(m courseFeedMeeting) string {
	switch {
	case m.CampusLocation == "O" && m.PmCode == "":
		return "Asynchronous Content :computer: Online"
	case m.CampusLocation == "O" && m.PmCode != "":
		return fmt.Sprintf("%s %s - %s :computer: Online",
			dayMap[m.MeetingDay], formatTime(m.StartTimeMilitary), formatTime(m.EndTimeMilitary))
	case m.CampusLocation == "3" && m.PmCode == "":
		return fmt.Sprintf("Hours by arrangement %s %s",
			campusMap[m.CampusName], m.CampusName)
	case m.CampusLocation != "" && m.PmCode != "":
		location := m.BuildingCode + "-" + m.RoomNumber
		if m.CampusName == "** INVALID **" {
			location = "Unknown"
		}
		return fmt.Sprintf("%s %s - %s %s **%s** (`%s`)",
			dayMap[m.MeetingDay], formatTime(m.StartTimeMilitary), formatTime(m.EndTimeMilitary),
			campusMap[m.CampusName], m.CampusName, location)
	default:
		return "Hours by arrangement"
	}
}

// todo: move this to discord search command handler once more interfaces are added
func formatTime(military string) string {
	rawHour, _ := strconv.Atoi(military[:2])
	hour := 12
	if rawHour%12 != 0 {
		hour = rawHour % 12
	}
	suffix := "AM"
	if rawHour/12 >= 1 {
		suffix = "PM"
	}
	return strconv.Itoa(hour) + ":" + military[2:] + suffix
}
