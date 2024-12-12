package shared

type Course struct {
	Title        string `json:"title"`
	CourseString string `json:"courseString"`
	Sections     []struct {
		Index        string `json:"index"`
		Section      string `json:"number"`
		Instructors  string `json:"instructorsText"`
		SectionNotes string `json:"sectionNotes"`
		OpenStatus   string   `json:"openStatus"`
		MeetingTimes []struct {
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

type CourseEntry struct {
	Title        string `json:"title"`
	CourseString string `json:"courseString"`
	Index        string `json:"index"`
	Section      string `json:"section"`
	Instructors  string `json:"instructors"`
	Notes        string `json:"notes"`
	Meeting      string `json:"meeting"`
}
