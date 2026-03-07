package user

import (
	"tsniper/internal/domain/scope"
)

type User struct {
	ID     string
	Campus scope.Campus
}

func NewUser(id string) *User {
	return &User{
		ID: id,
	}
}

func (u *User) DefaultCampus() *scope.Campus {
	if u.Campus == "" {
		return nil
	}
	return &u.Campus
}

func (u *User) SetCampus(campus scope.Campus) {
	u.Campus = campus
}

func (u *User) ClearCampus() {
	u.Campus = ""
}
