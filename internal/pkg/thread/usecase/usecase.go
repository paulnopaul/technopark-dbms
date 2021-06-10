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
func (t *threadUsecase) CreatePosts(s domain.SlugOrId, posts []domain.Post) ([]domain.Post, error) {
	panic("implement me")
}

func (t *threadUsecase) GetThreadDetails(s domain.SlugOrId, useSlug bool) (*domain.Thread, error) {
	panic("implement me")
}

func (t *threadUsecase) UpdateThreadDetails(s domain.SlugOrId, threadUpdate domain.Thread) (*domain.Thread, error) {
	panic("implement me")
}

func (t *threadUsecase) GetThreadPosts(s domain.SlugOrId, params domain.ArrayOutParams) {
	panic("implement me")
}

func (t *threadUsecase) VoteThread(s domain.SlugOrId, vote domain.Vote) {
	panic("implement me")
}
