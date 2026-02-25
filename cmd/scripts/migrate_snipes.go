package main

import (
	"database/sql"
	"fmt"

	"tsniper/internal/infrastructure/persistence/sqlite"
	"tsniper/pkg/database"
	"tsniper/pkg/logger"

	_ "modernc.org/sqlite"
)

const (
	summerTerm  = "7"
	currentYear = "2026"
)

var seasonToTerm = map[string]string{
	"summer": summerTerm,
	"spring": "1",
	"fall":   "9",
	"winter": "0",
}

type oldSnipe struct {
	CourseIndex string
	UserID      string
	Campus      string
	Season      string
}

func migrateSnipes() {
	log := logger.NewStdLogger(logger.LevelInfo)

	oldDB, err := sql.Open("sqlite", "./database_old.db")
	if err != nil {
		log.Fatal("Failed to open old database: %v", err)
	}
	defer oldDB.Close()

	newDB, err := database.NewSqliteDB("./database.db")
	if err != nil {
		log.Fatal("Failed to open new database: %v", err)
	}
	defer newDB.Close()

	if err := sqlite.Migrate(newDB); err != nil {
		log.Fatal("Failed to run migrations: %v", err)
	}

	rows, err := oldDB.Query("SELECT course_index, user_id, campus, season FROM snipes")
	if err != nil {
		log.Fatal("Failed to query old snipes: %v", err)
	}
	defer rows.Close()

	var snipes []oldSnipe
	for rows.Next() {
		var s oldSnipe
		if err := rows.Scan(&s.CourseIndex, &s.UserID, &s.Campus, &s.Season); err != nil {
			log.Fatal("Failed to scan snipe: %v", err)
		}
		snipes = append(snipes, s)
	}
	if err := rows.Err(); err != nil {
		log.Fatal("Row iteration error: %v", err)
	}

	log.Info("Found %d snipes in old database", len(snipes))

	tx, err := newDB.Begin()
	if err != nil {
		log.Fatal("Failed to begin transaction: %v", err)
	}

	stmt, err := tx.Prepare(
		`INSERT OR IGNORE INTO snipes (user_id, course_index, campus, term, year)
		 VALUES (?, ?, ?, ?, ?)`,
	)
	if err != nil {
		tx.Rollback()
		log.Fatal("Failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	migrated := 0
	skipped := 0
	for _, s := range snipes {
		term, ok := seasonToTerm[s.Season]
		if !ok {
			log.Warn("Unknown season %q for snipe (user=%s, index=%s), skipping", s.Season, s.UserID, s.CourseIndex)
			skipped++
			continue
		}

		result, err := stmt.Exec(s.UserID, s.CourseIndex, s.Campus, term, currentYear)
		if err != nil {
			log.Warn("Failed to insert snipe (user=%s, index=%s): %v", s.UserID, s.CourseIndex, err)
			skipped++
			continue
		}

		affected, _ := result.RowsAffected()
		if affected == 0 {
			log.Warn("Duplicate snipe skipped (user=%s, index=%s)", s.UserID, s.CourseIndex)
			skipped++
		} else {
			migrated++
		}
	}

	if err := tx.Commit(); err != nil {
		log.Fatal("Failed to commit transaction: %v", err)
	}

	log.Info("Snipe migration complete: %d migrated, %d skipped (out of %d total)",
		migrated, skipped, len(snipes))

	row := oldDB.QueryRow("SELECT COUNT(*) FROM snipes")
	var oldCount int
	row.Scan(&oldCount)

	newRow, _ := newDB.QueryRow("SELECT COUNT(*) FROM snipes")
	var newCount int
	newRow.Scan(&newCount)

	fmt.Printf("Verification: old=%d, new=%d\n", oldCount, newCount)
}
