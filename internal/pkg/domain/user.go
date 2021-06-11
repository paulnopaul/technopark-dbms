package domain

type User struct {
	Nickname string
	Fullname string
	About    string
	Email    string
}

type UserUsecase interface {
	CreateUser(nickname string, createData User) (*User, error)
	GetProfile(nickname string) (*User, error)
	UpdateProfile(nickname string, profileUpdate User) (*User, error)
}
