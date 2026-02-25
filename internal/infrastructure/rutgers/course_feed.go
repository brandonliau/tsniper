package rutgers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"tsniper/internal/application/ports"
	"tsniper/internal/domain/course"
	"tsniper/internal/domain/scope"

	"tsniper/pkg/httpx"
)

var _ ports.CourseFeed = (*courseFeed)(nil)

const coursesURL = "https://sis.rutgers.edu/soc/api/courses.json"

type courseFeedData struct {
	Title        string `json:"title"`
	CourseString string `json:"courseString"`
	Sections     []struct {
		Index        string              `json:"index"`
		Section      string              `json:"number"`
		Instructors  string              `json:"instructorsText"`
		SectionNotes string              `json:"sectionNotes"`
		OpenStatus   bool                `json:"openStatus"`
		MeetingTimes []courseFeedMeeting `json:"meetingTimes"`
	} `json:"sections"`
}

type courseFeedMeeting struct {
	CampusLocation    string `json:"campusLocation"`
	CampusName        string `json:"campusName"`
	PmCode            string `json:"pmCode"`
	MeetingDay        string `json:"meetingDay"`
	BuildingCode      string `json:"buildingCode"`
	RoomNumber        string `json:"roomNumber"`
	StartTimeMilitary string `json:"startTimeMilitary"`
	EndTimeMilitary   string `json:"endTimeMilitary"`
}

type courseFeed struct {
	client httpx.Client
}

func NewCourseFeed() *courseFeed {
	return &courseFeed{
		client: httpx.NewRetryClient(),
	}
}

func (f *courseFeed) FetchCourses(scp scope.AcademicScope) ([]*course.Course, error) {
	params := url.Values{}
	params.Add("campus", string(scp.Campus))
	params.Add("term", string(scp.Term))
	params.Add("year", scp.Year)

	fullURL := fmt.Sprintf("%s?%s", coursesURL, params.Encode())
	resp, err := f.client.Get(fullURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("course feed returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var courseFeed []courseFeedData
	if err := json.Unmarshal(body, &courseFeed); err != nil {
		return nil, err
	}

	return courseFeedDataToCourses(courseFeed, scp), nil
}
