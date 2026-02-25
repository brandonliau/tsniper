package sqlite

import (
	"database/sql"
	"errors"

	"tsniper/internal/domain/user"

	"tsniper/pkg/database"
)

var _ user.UserRepository = (*UserRepositoryImpl)(nil)

type UserRepositoryImpl struct {
	db *database.SqliteDB
}

func NewUserRepository(db *database.SqliteDB) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		db: db,
	}
}

func (r *UserRepositoryImpl) Create(usr *user.User) error {
	return r.db.Exec(
		`INSERT OR IGNORE INTO users (user_id, campus)
		 VALUES (?, ?)`,
		usr.ID,
		usr.Campus,
	)
}

func (r *UserRepositoryImpl) Save(usr *user.User) error {
	return r.db.Exec(
		`UPDATE users
		 SET campus = ?
		 WHERE user_id = ?`,
		usr.Campus,
		usr.ID,
	)
}

func (r *UserRepositoryImpl) Delete(usr *user.User) error {
	return r.db.Exec(
		`DELETE FROM users
		 WHERE user_id = ?`,
		usr.ID,
	)
}

func (r *UserRepositoryImpl) Get(userID string) (*user.User, error) {
	row, err := r.db.QueryRow(
		`SELECT user_id, campus
		 FROM users
		 WHERE user_id = ?`,
		userID,
	)
	if err != nil {
		return nil, err
	}

	var usr user.User
	err = row.Scan(&usr.ID, &usr.Campus)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, user.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return &usr, nil
}

func (r *UserRepositoryImpl) GetAll() ([]*user.User, error) {
	rows, err := r.db.Query(
		`SELECT user_id, campus
		 FROM users`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*user.User
	for rows.Next() {
		var usr user.User
		err = rows.Scan(&usr.ID, &usr.Campus)
		if err != nil {
			return nil, err
		}
		users = append(users, &usr)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
