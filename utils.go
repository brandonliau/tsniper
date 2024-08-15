package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
)

/* API ENDPOINTS */
func OpenSectionsAPI(campus string, season Season, client *http.Client) []string {
	var openSections []string
	defer func() {
		if r := recover(); r != nil {
			openSections = make([]string, 0)
			fmt.Printf("Recovered @ %s : %s\n", time.Now().Format("2006-01-02 15:04:05.00000"), r)
		}
	}()
	var url string = fmt.Sprintf(
		"https://sis.rutgers.edu/soc/api/openSections.json?year=%s&term=%s&campus=%s",
		season.Year,
		season.Term,
		campus,
	)
	resp, err := client.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	rawJson, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(rawJson, &openSections)
	return openSections
}

func CoursesAPI(campus string, season Season) []RawCourseData {
	var rawCourseData []RawCourseData
	defer func() {
		if r := recover(); r != nil {
			rawCourseData = make([]RawCourseData, 0)
			fmt.Printf("Recovered @ %s : %s\n", time.Now().Format("2006-01-02 15:04:05.00000"), r)
		}
	}()
	var url string = fmt.Sprintf(
		"https://sis.rutgers.edu/soc/api/courses.json?year=%s&term=%s&campus=%s",
		season.Year,
		season.Term,
		campus,
	)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	rawJson, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(rawJson, &rawCourseData)
	return rawCourseData
}

/* SUPPORT FUNCTIONS */
func (runTimeData *RunTimeData) LoadConfig(filename string) {
	rawYaml, _ := os.ReadFile(filename)
	_ = yaml.Unmarshal(rawYaml, &runTimeData.Config)
	fmt.Printf("SUCCESS @ %s : LOAD CONFIG\n", time.Now().Format("2006-01-02 15:04:05.00000"))
}

func (runTimeData *RunTimeData) InitRunTimeData() {
	// open connection to db
	db, _ := sql.Open("sqlite", "./database.db")
	runTimeData.Db = db
	fmt.Printf("SUCCESS @ %s : CONNECTED TO DATABASE\n", time.Now().Format("2006-01-02 15:04:05.00000"))
	// intialize database
	runTimeData.InitDb()
	fmt.Printf("SUCCESS @ %s : INITIALIZE DATABASE\n", time.Now().Format("2006-01-02 15:04:05.00000"))
	// intialize Tracking and PrevOpened maps
	runTimeData.Tracking = make(map[string]map[string]int)
}

func (runTimeData *RunTimeData) SyncUsers() {
	rows := runTimeData.GetDistinctUsers()
	var memberID string
	for rows.Next() {
		rows.Scan(&memberID)
		_, err := s.State.Member(runTimeData.Config.Guild, memberID)
		if err != nil {
			runTimeData.ClearSnipe(memberID)
		}
	}
	fmt.Printf("SUCCESS @ %s : SYNC USERS\n", time.Now().Format("2006-01-02 15:04:05.00000"))
}

func (runTimeData *RunTimeData) UpdateTracking(action int, index string, campus string, season string) {
	switch action {
	case 0: // sync
		for _, campus := range runTimeData.Config.CurrentCampuses {
			for _, season := range runTimeData.Config.CurrentSeasons {
				runTimeData.Tracking[campus + season] = make(map[string]int)
			}
		}
		rows := runTimeData.GetDistinctSnipes()
		var index, campus, season string
		var count int
		for rows.Next() {
			rows.Scan(&index, &campus, &season, &count);
			runTimeData.Tracking[campus + season][index] = count
		}
	case 1: // add
		runTimeData.Tracking[campus + season][index] += 1
	case 2: // remove
		runTimeData.Tracking[campus + season][index] -= 1
		if runTimeData.Tracking[campus + season][index] == 0 {
			delete(runTimeData.Tracking[campus + season], index)
		}
	case 3: // remove key (course open event)
		delete(runTimeData.Tracking[campus + season], index)
	}
}

func (runTimeData *RunTimeData) SyncTracking() {
	for _, campus := range runTimeData.Config.CurrentCampuses {
		for _, season := range runTimeData.Config.CurrentSeasons {
			runTimeData.Tracking[campus + season] = make(map[string]int)
			runTimeData.UpdateTracking(0, "", campus, season)
		}
	}
}

func (runTimeData *RunTimeData) GetCampus(memberID string) string {
	member, _ := s.State.Member(runTimeData.Config.Guild, memberID)
	roles := member.Roles
	for _, roleId := range roles {
		role, _ := s.State.Role(runTimeData.Config.Guild, roleId)
		switch role.Name {
		case "New-Brunswick":
			return "NB"
		case "Newark":
			return "NK"
		case "Camden":
			return "CM"
		}
	}
	return "NB"
}

func (runTimeData *RunTimeData) GetSeason(season string) string {
	for k, v := range runTimeData.Config.Seasons {
		if v.Term + v.Year == season {
			return k
		}
	}
	return ""
}

func ParseTime(time string) string {
	rawHour, _ := strconv.Atoi(time[:2])
	var hour int = 12
	if rawHour % 12 != 0 {
		hour = rawHour % 12
	}
	timeString := strconv.Itoa(hour) + ":" + time[2:]
    if rawHour / 12 < 1 {
        timeString += "AM"
	} else {
        timeString += "PM"
	}
    return timeString
}

func GetKeys[k comparable, v any](data map[k]v) []k {
	keys := make([]k, len(data))
	i := 0
	for key := range data {
		keys[i] = key
		i++ 
	}
	return keys
}

func Intersection[T comparable](data1 []T, data2 []T) []T {
	intersection := make([]T, 0)
	hash := make(map[T]struct{})
	for _, i := range data1 {
		hash[i] = struct{}{}
	}
	for _, j := range data2 {
		if _, ok := hash[j]; ok {
			intersection = append(intersection, j)
		}
	}
	return intersection
}

// Returns the elements in `data1` that aren't in `data2`
func Difference[T comparable](data1 []T, data2 []T) []T {
	difference := make([]T, 0)
	hash := make(map[T]struct{}, len(data2))
	for _, i := range data2 {
		hash[i] = struct{}{}
	}
	for _, j := range data1 {
		if _, ok := hash[j]; !ok {
			difference = append(difference, j)
		}
	}
	return difference
}