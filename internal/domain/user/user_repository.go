package user

type UserRepository interface {
	Create(usr *User) error
	Save(usr *User) error
	Delete(usr *User) error

	Get(userID string) (*User, error)
	GetAll() ([]*User, error)
}
