package usecase

import (
	"DBMSForum/internal/pkg/domain"
	"database/sql"
)

type userUsecase struct {
	DB *sql.DB
}

func NewUserUsecase(db *sql.DB) domain.UserUsecase {
	return &userUsecase{
		DB: db,
	}
}

func (u *userUsecase) CreateUser(nickname string, createData domain.User) (*domain.User, error) {
	panic("implement me")
}

func (u *userUsecase) GetProfile(nickname string) (*domain.User, error) {
	panic("implement me")
}

func (u *userUsecase) UpdateProfile(nickname string, profileUpdate domain.User) (*domain.User, error) {
	panic("implement me")
}
