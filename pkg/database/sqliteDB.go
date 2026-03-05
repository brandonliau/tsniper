package database

import (
	"database/sql"
	"errors"
	"os"
)

type SqliteDB struct {
	readDB  *sql.DB
	writeDB *sql.DB
}

func NewSqliteDB(file string) (*SqliteDB, error) {
	readDB, err := sql.Open("sqlite", file)
	if err != nil {
		return nil, err
	}
	writeDB, err := sql.Open("sqlite", file)
	if err != nil {
		return nil, err
	}
	sqliteDB := &SqliteDB{
		readDB:  readDB,
		writeDB: writeDB,
	}
	err = sqliteDB.applyPerformanceOptions(readDB, 4, "./pkg/database/sqlitePerf.sql")
	if err != nil {
		return nil, err
	}
	err = sqliteDB.applyPerformanceOptions(writeDB, 1, "./pkg/database/sqlitePerf.sql")
	if err != nil {
		return nil, err
	}
	return sqliteDB, nil
}

func (db *SqliteDB) applyPerformanceOptions(sdb *sql.DB, maxOpenConns int, file string) error {
	sdb.SetMaxOpenConns(maxOpenConns)
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	_, err = sdb.Exec(string(data))
	if err != nil {
		return err
	}
	return nil
}

func (db *SqliteDB) Close() error {
	readErr := db.readDB.Close()
	writeErr := db.writeDB.Close()
	return errors.Join(readErr, writeErr)
}

func (db *SqliteDB) Exec(query string, args ...any) error {
	_, err := db.writeDB.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (db *SqliteDB) ExecSQLFile(file string) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	_, err = db.writeDB.Exec(string(data))
	if err != nil {
		return err
	}
	return nil
}

func (db *SqliteDB) Query(query string, args ...any) (*sql.Rows, error) {
	rows, err := db.readDB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (db *SqliteDB) QueryRow(query string, args ...any) (*sql.Row, error) {
	row := db.readDB.QueryRow(query, args...)
	return row, nil
}

func (db *SqliteDB) PrepareExec(query string) (*sql.Stmt, error) {
	stmt, err := db.writeDB.Prepare(query)
	if err != nil {
		return nil, err
	}
	return stmt, nil
}

func (db *SqliteDB) PrepareQuery(query string) (*sql.Stmt, error) {
	stmt, err := db.readDB.Prepare(query)
	if err != nil {
		return nil, err
	}
	return stmt, nil
}

func (db *SqliteDB) Begin() (tx *sql.Tx, err error) {
	return db.writeDB.Begin()
}
