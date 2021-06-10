package usecase

import (
	"DBMSForum/internal/pkg/domain"
	"database/sql"
)

type threadUsecase struct {
	DB *sql.DB
}

func NewThreadUsecase(db *sql.DB) domain.ThreadUsecase {
	return &threadUsecase{
		DB: db,
	}
}
