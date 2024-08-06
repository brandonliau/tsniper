package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

func (runTimeData *RunTimeData) GetOpenSections(openSections chan []string, wg *sync.WaitGroup, campus string, season string, client *http.Client) {
	defer wg.Done()
	openData := OpenSectionsAPI(campus, runTimeData.Config.Seasons[season], client)
	trackingData := GetKeys(runTimeData.Tracking[campus + season])
	data := Intersection(openData, trackingData)
	if len(data) > 0 {
		data = append(data, campus + season)
		openSections <- data
	}
	runTimeData.mtx.Lock()
	defer runTimeData.mtx.Unlock()
	if len(runTimeData.PrevOpened[campus + season]) == 0 {
		runTimeData.PrevOpened[campus + season] = openData
	}
	openDiff := Difference(openData, runTimeData.PrevOpened[campus + season])
	runTimeData.PrevOpened[campus + season] = append(runTimeData.PrevOpened[campus + season], openDiff...)
	closeDiff := Difference(runTimeData.PrevOpened[campus + season], openData)
	for _, index := range closeDiff {
		runTimeData.UpdateLastOpen(index, campus, season)
		runTimeData.PrevOpened[campus + season] = openData
	}
}

func (runTimeData *RunTimeData) SnipeCheck() {
	transport := &http.Transport{
		MaxConnsPerHost: 6,
        MaxIdleConnsPerHost: 6,
    }
	client := &http.Client{Transport: transport}
	for {
		select {
		case <-snipeTicker.C:
			openSections := make(chan []string)
			var wg sync.WaitGroup
			for _, campus := range runTimeData.Config.CurrentCampuses {
				for _, season := range runTimeData.Config.CurrentSeasons {
					wg.Add(1)
					go runTimeData.GetOpenSections(openSections, &wg, campus, season, client)
				}
			}
			go func() {
				wg.Wait()
				close(openSections)
			}()
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
		case <-snipeClose:
			return
		}
	}
}
