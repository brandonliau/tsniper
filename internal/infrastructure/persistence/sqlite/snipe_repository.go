package sqlite

import (
	"database/sql"
	"errors"

	"tsniper/internal/domain/scope"
	"tsniper/internal/domain/snipe"

	"tsniper/pkg/database"

	"modernc.org/sqlite"
	sqlitelib "modernc.org/sqlite/lib"
)

var _ snipe.SnipeRepository = (*SnipeRepositoryImpl)(nil)

type SnipeRepositoryImpl struct {
	db *database.SqliteDB
}

func NewSnipeRepository(db *database.SqliteDB) *SnipeRepositoryImpl {
	return &SnipeRepositoryImpl{
		db: db,
	}
}

func (r *SnipeRepositoryImpl) Create(snp *snipe.Snipe) error {
	err := r.db.Exec(
		`INSERT INTO snipes (user_id, course_index, campus, term, year)
		 VALUES (?, ?, ?, ?, ?)`,
		snp.UserID,
		snp.Index,
		snp.Scope.Campus,
		snp.Scope.Term,
		snp.Scope.Year,
	)
	if err != nil {
		var sqliteErr *sqlite.Error
		if errors.As(err, &sqliteErr) && sqliteErr.Code() == sqlitelib.SQLITE_CONSTRAINT_UNIQUE {
			return snipe.ErrSnipeDuplicate
		}
		return err
	}

	return nil
}

func (r *SnipeRepositoryImpl) Delete(snp *snipe.Snipe) error {
	return r.db.Exec(
		`DELETE FROM snipes
		 WHERE user_id = ? AND course_index = ? AND campus = ? AND term = ? AND year = ?`,
		snp.UserID,
		snp.Index,
		snp.Scope.Campus,
		snp.Scope.Term,
		snp.Scope.Year,
	)
}

func (r *SnipeRepositoryImpl) Get(userID string, index string, scp scope.AcademicScope) (*snipe.Snipe, error) {
	row, err := r.db.QueryRow(
		`SELECT user_id, course_index, campus, term, year
		 FROM snipes
		 WHERE user_id = ? AND course_index = ? AND campus = ? AND term = ? AND year = ?`,
		userID,
		index,
		scp.Campus,
		scp.Term,
		scp.Year,
	)
	if err != nil {
		return nil, err
	}

	var snp snipe.Snipe
	err = row.Scan(&snp.UserID, &snp.Index, &snp.Scope.Campus, &snp.Scope.Term, &snp.Scope.Year)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, snipe.ErrSnipeNotFound
	}
	if err != nil {
		return nil, err
	}

	return &snp, nil
}

func (r *SnipeRepositoryImpl) ListAll() ([]*snipe.Snipe, error) {
	rows, err := r.db.Query(
		`SELECT user_id, course_index, campus, term, year
		 FROM snipes`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var snipes []*snipe.Snipe
	for rows.Next() {
		var snp snipe.Snipe
		err = rows.Scan(&snp.UserID, &snp.Index, &snp.Scope.Campus, &snp.Scope.Term, &snp.Scope.Year)
		if err != nil {
			return nil, err
		}
		snipes = append(snipes, &snp)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return snipes, nil
}

func (r *SnipeRepositoryImpl) ListByUser(userID string) ([]*snipe.Snipe, error) {
	rows, err := r.db.Query(
		`SELECT user_id, course_index, campus, term, year
		 FROM snipes
		 WHERE user_id = ?`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var snipes []*snipe.Snipe
	for rows.Next() {
		var snp snipe.Snipe
		err = rows.Scan(&snp.UserID, &snp.Index, &snp.Scope.Campus, &snp.Scope.Term, &snp.Scope.Year)
		if err != nil {
			return nil, err
		}
		snipes = append(snipes, &snp)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return snipes, nil
}

func (r *SnipeRepositoryImpl) ListByIndex(index string, scp scope.AcademicScope) ([]*snipe.Snipe, error) {
	rows, err := r.db.Query(
		`SELECT user_id, course_index, campus, term, year
		 FROM snipes
		 WHERE course_index = ? AND campus = ? AND term = ? AND year = ?`,
		index,
		scp.Campus,
		scp.Term,
		scp.Year,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var snipes []*snipe.Snipe
	for rows.Next() {
		var snp snipe.Snipe
		err = rows.Scan(&snp.UserID, &snp.Index, &snp.Scope.Campus, &snp.Scope.Term, &snp.Scope.Year)
		if err != nil {
			return nil, err
		}
		snipes = append(snipes, &snp)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return snipes, nil
}

func (r *SnipeRepositoryImpl) DeleteByUser(userID string) error {
	return r.db.Exec(
		`DELETE FROM snipes
		 WHERE user_id = ?`,
		userID,
	)
}
