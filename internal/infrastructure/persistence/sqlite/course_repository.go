package sqlite

import (
	"database/sql"
	"errors"

	"tsniper/internal/domain/course"
	"tsniper/internal/domain/scope"

	"tsniper/pkg/database"
)

var _ course.CourseRepository = (*CourseRepositoryImpl)(nil)

type CourseRepositoryImpl struct {
	db *database.SqliteDB
}

func NewCourseRepository(db *database.SqliteDB) *CourseRepositoryImpl {
	return &CourseRepositoryImpl{
		db: db,
	}
}

func (r *CourseRepositoryImpl) Create(crs *course.Course) error {
	return r.db.Exec(
		`INSERT OR IGNORE INTO courses (course_index, course_string, section, title, instructors, notes, meeting, campus, term, year, last_open)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		crs.Index,
		crs.CourseString,
		crs.Section,
		crs.Title,
		crs.Instructors,
		crs.Notes,
		crs.Meeting,
		crs.Scope.Campus,
		crs.Scope.Term,
		crs.Scope.Year,
		crs.LastOpen,
	)
}

func (r *CourseRepositoryImpl) Save(crs *course.Course) error {
	return r.db.Exec(
		`UPDATE courses
		 SET course_string = ?, section = ?, title = ?, instructors = ?, notes = ?, meeting = ?, last_open = ?
		 WHERE course_index = ? AND campus = ? AND term = ? AND year = ?`,
		crs.CourseString,
		crs.Section,
		crs.Title,
		crs.Instructors,
		crs.Notes,
		crs.Meeting,
		crs.LastOpen,
		crs.Index,
		crs.Scope.Campus,
		crs.Scope.Term,
		crs.Scope.Year,
	)
}

func (r *CourseRepositoryImpl) Delete(course *course.Course) error {
	return r.db.Exec(
		`DELETE FROM courses
		 WHERE course_index = ? AND campus = ? AND term = ? AND year = ?`,
		course.Index,
		course.Scope.Campus,
		course.Scope.Term,
		course.Scope.Year,
	)
}

func (r *CourseRepositoryImpl) Get(index string, scp scope.AcademicScope) (*course.Course, error) {
	row, err := r.db.QueryRow(
		`SELECT course_index, course_string, section, title, instructors, notes, meeting, campus, term, year, last_open
		 FROM courses
		 WHERE course_index = ? AND campus = ? AND term = ? AND year = ?`,
		index,
		scp.Campus,
		scp.Term,
		scp.Year,
	)
	if err != nil {
		return nil, err
	}

	var crs course.Course
	err = row.Scan(&crs.Index, &crs.CourseString, &crs.Section, &crs.Title, &crs.Instructors, &crs.Notes, &crs.Meeting, &crs.Scope.Campus, &crs.Scope.Term, &crs.Scope.Year, &crs.LastOpen)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, course.ErrCourseNotFound
	}
	if err != nil {
		return nil, err
	}

	return &crs, nil
}

func (r *CourseRepositoryImpl) GetAll() ([]*course.Course, error) {
	rows, err := r.db.Query(
		`SELECT course_index, course_string, section, title, instructors, notes, meeting, campus, term, year, last_open
		 FROM courses`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []*course.Course
	for rows.Next() {
		var crs course.Course
		err = rows.Scan(&crs.Index, &crs.CourseString, &crs.Section, &crs.Title, &crs.Instructors, &crs.Notes, &crs.Meeting, &crs.Scope.Campus, &crs.Scope.Term, &crs.Scope.Year, &crs.LastOpen)
		if err != nil {
			return nil, err
		}
		courses = append(courses, &crs)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return courses, nil
}

func (r *CourseRepositoryImpl) GetAllByScope(scp scope.AcademicScope) ([]*course.Course, error) {
	rows, err := r.db.Query(
		`SELECT course_index, course_string, section, title, instructors, notes, meeting, campus, term, year, last_open
		 FROM courses
		 WHERE campus = ? AND term = ? AND year = ?`,
		scp.Campus,
		scp.Term,
		scp.Year,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []*course.Course
	for rows.Next() {
		var crs course.Course
		err = rows.Scan(&crs.Index, &crs.CourseString, &crs.Section, &crs.Title, &crs.Instructors, &crs.Notes, &crs.Meeting, &crs.Scope.Campus, &crs.Scope.Term, &crs.Scope.Year, &crs.LastOpen)
		if err != nil {
			return nil, err
		}
		courses = append(courses, &crs)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return courses, nil
}

func (r *CourseRepositoryImpl) BatchCreate(courses []*course.Course) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(
		`INSERT OR IGNORE INTO courses (course_index, course_string, section, title, instructors, notes, meeting, campus, term, year, last_open)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, crs := range courses {
		if _, err := stmt.Exec(crs.Index, crs.CourseString, crs.Section, crs.Title, crs.Instructors, crs.Notes, crs.Meeting, crs.Scope.Campus, crs.Scope.Term, crs.Scope.Year, crs.LastOpen); err != nil {
			return err
		}
	}

	return tx.Commit()
}
