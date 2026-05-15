package view

import (
	"tsniper/internal/domain/user"
)

type UserView struct {
	ID     string
	Campus string
}

func FromUser(u *user.User) *UserView {
	return &UserView{
		ID:     u.ID,
		Campus: string(u.Campus),
	}
}

func FromUsers(users []*user.User) []*UserView {
	out := make([]*UserView, 0, len(users))
	for _, u := range users {
		out = append(out, FromUser(u))
	}
	return out
}
