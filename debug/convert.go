package debug

import (
	"database/sql"
	"strconv"
	
	_ "modernc.org/sqlite"
)

type InputData struct {
	Indices []string	`json:"indices"`
	Ruid	string		`json:"ruid"`
	Pac		string		`json:"pac"`
}

type OutputData struct {
	Indices map[string]struct{}
	Ruid	string
	Pac		string
}

var seasonMap = map[string]string{
	"92024": "FALL",
	"72024": "SUMMER",
}

func convert() {
	old_db, _ := sql.Open("sqlite", "./old_database.db")
	new_db, _ := sql.Open("sqlite", "./new_database.db")
	rows, _ := old_db.Query("SELECT course_index, user_id, campus, term, year FROM snipes")
	var course_index, campus string
	var user_id, term, year int
	for rows.Next() {
		rows.Scan(&course_index, &user_id, &campus, &term, &year)
		new_user_id := strconv.Itoa(user_id)
		new_term := strconv.Itoa(term)
		new_year := strconv.Itoa(year)
		season := seasonMap[new_term + new_year]
		new_db.Exec("INSERT INTO snipes (course_index, memberID, campus, season) VALUES (?, ?, ?, ?)",
			course_index,
			new_user_id,
			campus,
			season,
		)
	}
}

func main() {
	convert()
}