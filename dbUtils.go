package main

import (
	"fmt"
	"sort"
	"time"

	"github.com/bwmarrin/discordgo"
)

func (runTimeData *RunTimeData) InitDb() {
	snipesDb := "CREATE TABLE IF NOT EXISTS snipes (course_index TEXT, memberID TEXT, campus TEXT, season TEXT)"
	runTimeData.Db.Exec(snipesDb)

	registeredDb := "CREATE TABLE IF NOT EXISTS commands (command_name TEXT, command_id TEXT)"
	runTimeData.Db.Exec(registeredDb)
	
	for _, campus := range runTimeData.Config.CurrentCampuses {
		campusDb := fmt.Sprintf(
			"CREATE TABLE IF NOT EXISTS %s" + 
			"(course_index TEXT, title TEXT, course_string TEXT, section TEXT, instructors TEXT, notes TEXT, meeting TEXT, season TEXT, last_open BIGINT, " + 
			"UNIQUE (course_index, season))",
			campus)
		runTimeData.Db.Exec(campusDb)
	}
}

func (runTimeData *RunTimeData) AddSnipe(index string, memberID string, campus string, season string) {
	query := "INSERT INTO snipes (course_index, memberID, campus, season) VALUES (?, ?, ?, ?)"
	runTimeData.Db.Exec(query, index, memberID, campus, season)
}

func (runTimeData *RunTimeData) RemoveSnipe(index string, memberID string, season string) {
	query := "DELETE FROM snipes WHERE course_index = ? AND memberID = ? AND season = ?"
	runTimeData.Db.Exec(query, index, memberID, season)
}

func (runTimeData *RunTimeData) ClearSnipe(memberID string) {
	query := "DELETE FROM snipes WHERE memberID = ?"
	runTimeData.Db.Exec(query, memberID)
}

func (runTimeData *RunTimeData) GetSnipes(memberID string) [][]string {
	var index, campus, season string
	snipes := make([][]string, 0)
	query := "SELECT course_index, campus, season FROM snipes WHERE memberID = ?"
	rows, _ := runTimeData.Db.Query(query, memberID)
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&index, &campus, &season)
		snipes = append(snipes, []string{index, campus, season})
	}
	sort.SliceStable(snipes, func(i, j int) bool {
		return snipes[i][0] < snipes[j][0]
	})
	return snipes
}

func (runTimeData *RunTimeData) CheckCourseExist(index string, campus string, season string) bool {
	var exist int
	query := fmt.Sprintf("SELECT 1 FROM %s WHERE course_index = ? AND season = ?", campus)
	row := runTimeData.Db.QueryRow(query, index, season)
	err := row.Scan(&exist)
	return err == nil // return true if course exists in db
}

func (runTimeData *RunTimeData) GetCourseData(index string, campus string, season string) CourseData {
	var course_index, title, courseString, section, instructors, notes, meeting string
	query := fmt.Sprintf("SELECT course_index, title, course_string, section, instructors, notes, meeting FROM %s WHERE course_index = ? AND season = ?", campus)
	row := runTimeData.Db.QueryRow(query, index, season)
	row.Scan(&course_index, &title, &courseString, &section, &instructors, &notes, &meeting)
	return CourseData{
		Title: title, 
		CourseString: courseString,
		Index: course_index,
		Section: section,
		Instructors: instructors,
		Notes: notes,
		Meeting: meeting,
	}
}

func (runTimeData *RunTimeData) GetSnipeCount(index string, campus string, season string) int {
	var count int
	query := "SELECT count(*) FROM snipes WHERE course_index = ? AND campus = ? AND season = ?"
	row := runTimeData.Db.QueryRow(query, index, campus, season)
	row.Scan(&count)
	return count
}

func (runTimeData *RunTimeData) GetLastOpen(index string, campus string, season string) int64 {
	var lastOpen int64
	query := fmt.Sprintf("SELECT last_open FROM %s WHERE course_index = ? AND season = ?", campus)
	row := runTimeData.Db.QueryRow(query, index, season)
	row.Scan(&lastOpen)
	return lastOpen
}

func (runTimeData *RunTimeData) GetLastOpens(campus string, season string) map[string]int64 {
	var course_index string
	var lastOpen int64
	lastOpens := make(map[string]int64)
	query := fmt.Sprintf("SELECT course_index, last_open FROM %s WHERE season = ?", campus)
	rows, _ := runTimeData.Db.Query(query, season)
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&course_index, &lastOpen)
		lastOpens[course_index] = lastOpen
	}
	return lastOpens
}

func (runTimeData *RunTimeData) UpdateLastOpen(time int64, index string, campus string, season string) {
	query := fmt.Sprintf("UPDATE %s SET last_open = ? WHERE course_index = ? AND season = ?", campus)
	runTimeData.Db.Exec(query, time, index, season)
}

func (runTimeData *RunTimeData) GetUsersFromIndex(index string, campus string, season string) []string {
	var users []string
	var memberID string
	query := "SELECT memberID FROM snipes WHERE course_index = ? AND season = ?"
	rows, err := runTimeData.Db.Query(query, index, season)
	if err != nil {
		fmt.Println(err, time.Now().Format("2006-01-02 15:04:05.00000"))
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&memberID)
		users = append(users, memberID)
	}
	return users
}

func (runTimeData *RunTimeData) GetDistinctUsers() []string {
	var users []string
	var memberID string
	query := "SELECT DISTINCT memberID FROM snipes"
	rows, _ := runTimeData.Db.Query(query)
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&memberID)
		users = append(users, memberID)
	}
	return users
}

func (runTimeData *RunTimeData) GetDistinctSnipes() map[string]map[string]int {
	var index, campus, season string
	var count int
	tracking := make(map[string]map[string]int)
	query := "SELECT course_index, campus, season, COUNT(course_index) AS count FROM snipes GROUP BY course_index"
	rows, _ := runTimeData.Db.Query(query)
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&index, &campus, &season, &count);
		if tracking[campus + season] == nil {
			tracking[campus + season] = make(map[string]int)
		}
		tracking[campus + season][index] = count
	}
	return tracking
}

func (runTimeData *RunTimeData) AreadySniping(memberID string, index string, campus string, season string) bool {
	var exist int
	query := "SELECT 1 FROM snipes WHERE memberID = ? AND course_index = ? AND campus = ? AND season = ?"
	row := runTimeData.Db.QueryRow(query, memberID, index, campus, season)
	err := row.Scan(&exist)
	return err == nil // return true if snipe exists in db
}

func (runTimeData *RunTimeData) UpdateRegisteredCommands(registered []*discordgo.ApplicationCommand) {
	runTimeData.Db.Exec("DELETE FROM commands")
	for _, command := range registered {
		query := "INSERT INTO commands (command_name, command_id) VALUES (?, ?)"
		runTimeData.Db.Exec(query, command.Name, command.ID)
	}
}

func (runTimeData *RunTimeData) GetRegisteredCommands() map[string]string {
	var name, id string
	registered := make(map[string]string)
	query := "SELECT command_name, command_id FROM commands"
	rows, _ := runTimeData.Db.Query(query)
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&name, &id)
		registered[name] = id
	}
	return registered
}
