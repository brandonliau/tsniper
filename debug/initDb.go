package debug

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

var currentCampuses = []string{"NB", "NK", "CM"}

func InitDb() {
	db, _ := sql.Open("sqlite", "./database.db")
	snipeDb := "CREATE TABLE IF NOT EXISTS snipes (course_index TEXT, memberID TEXT, campus TEXT, season TEXT)"
	db.Prepare(snipeDb)
	for _, campus := range currentCampuses {
		campusQuery := fmt.Sprintf(
			"CREATE TABLE IF NOT EXISTS %s (course_index TEXT, title TEXT, course_string TEXT, section TEXT, instructors TEXT, notes TEXT, meeting TEXT, season TEXT, last_open BIGINT)",
			campus)
		db.Prepare(campusQuery)
	}
}
