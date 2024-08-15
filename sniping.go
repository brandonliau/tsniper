package main

import (
	"fmt"
	"net/http"
	"time"
)

func (runTimeData *RunTimeData) GetOpenSections(openSections chan []string, campus string, season string, client *http.Client) {
	prevOpened := OpenSectionsAPI(campus, runTimeData.Config.Seasons[season], client)
	for {
		select {
		case <- snipeTicker.C:
			// fmt.Println(campus, season, time.Now().Format("2006-01-02 15:04:05.00000"))
			openData := OpenSectionsAPI(campus, runTimeData.Config.Seasons[season], client)
			trackingData := GetKeys(runTimeData.Tracking[campus + season])
			openCourses := Intersection(openData, trackingData)
			if len(openCourses) > 0 {
				openCourses = append(openCourses, campus + season) // append [campus + season] for processing loop
				openSections <- openCourses
			}
			openDiff := Difference(openData, prevOpened)
			closeDiff := Difference(prevOpened, openData)
			for _, index := range openDiff {
				fmt.Println(index)
				runTimeData.UpdateLastOpen(0, index, campus, season)
			}
			for _, index := range closeDiff {
				fmt.Println(index)
				runTimeData.UpdateLastOpen(time.Now().Unix(), index, campus, season)
			}
			prevOpened = openData
		case <-snipeClose:
			close(openSections)
			return
		}
	}
}

func (runTimeData *RunTimeData) SnipeCheck() {
	transport := &http.Transport{
		MaxConnsPerHost: 6,
        MaxIdleConnsPerHost: 6,
    }
	client := &http.Client{Transport: transport}
	openSections := make(chan []string)
	for _, campus := range runTimeData.Config.CurrentCampuses {
		for _, season := range runTimeData.Config.CurrentSeasons {
			go runTimeData.GetOpenSections(openSections, campus, season, client)
		}
	}
	for data := range openSections {
		campusSeason := data[len(data) - 1]
		campus := campusSeason[0:2]
		season := campusSeason[2:]
		data = data[:len(data) - 1]
		for _, index := range data {
			course := runTimeData.GetCourseData(index, campus, season)
			embed := runTimeData.SnipeEmbed(course)
			regButton := runTimeData.CreateRegButton(course, runTimeData.Config.Seasons[season], campus)
			resnipeButton := CreateResnipeButton()
			for _, memberID := range runTimeData.GetUsersFromIndex(index, campus, season) {
				user, _ := s.User(memberID)
				dmChannel, _ := s.UserChannelCreate(memberID)
				err := SendComplexMessage(user, dmChannel, embed, regButton, resnipeButton)
				if err != nil {
					fmt.Printf("FAILURE @ %s : USER %s HASN'T ENABLED DIRECT MESSAGES\n", time.Now().Format("2006-01-02 15:04:05.00000"), memberID)
				}
				runTimeData.RemoveSnipe(index, memberID, season)
			}
			runTimeData.UpdateTracking(3, index, campus, season)
		}
	}
}
