package sqlite

import (
	"embed"
	"io/fs"

	"tsniper/pkg/database"
)

//go:embed migrations/*.sql
var migrationFS embed.FS

func Migrate(db *database.SqliteDB) error {
	entries, err := fs.ReadDir(migrationFS, "migrations")
	if err != nil {
		return err
	}

	for _, entry := range entries {
		data, err := fs.ReadFile(migrationFS, "migrations/"+entry.Name())
		if err != nil {
			return err
		}
		if err := db.Exec(string(data)); err != nil {
			return err
		}
	}

	return nil
}
