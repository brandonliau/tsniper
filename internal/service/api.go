package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"Tsniper/internal/shared"
)

func fetchData(client *http.Client, baseURL string, params url.Values, result any) error {
	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
	resp, err := client.Get(fullURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, result)
	if err != nil {
		return err
	}
	return nil
}

func Courses(client *http.Client, year, term, campus string) ([]shared.Course, error) {
	baseURL := "https://sis.rutgers.edu/soc/api/courses.json"
	params := url.Values{}
	params.Add("year", year)
	params.Add("term", term)
	params.Add("campus", campus)
	var courses []shared.Course
	err := fetchData(client, baseURL, params, &courses)
	if err != nil {
		return nil, err
	}
	return courses, nil
}

func OpenSections(client *http.Client, year, term, campus string) []string {
	baseURL := "https://sis.rutgers.edu/soc/api/openSections.json"
	params := url.Values{}
	params.Add("year", year)
	params.Add("term", term)
	params.Add("campus", campus)
	var openSections []string
	fetchData(client, baseURL, params, &openSections)
	return openSections
}
