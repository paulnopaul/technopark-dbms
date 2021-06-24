package usecase

import (
	"database/sql"
	"technopark-dbms/internal/pkg/domain"
	"technopark-dbms/internal/pkg/utilities"
)

type threadUsecase struct {
	DB *sql.DB
}

func (t threadUsecase) CreatePosts(s utilities.SlugOrId, posts []domain.Post) ([]domain.Post, error) {
	panic("implement me")
}

func (t threadUsecase) GetThreadDetails(s utilities.SlugOrId) (*domain.Thread, error) {
	panic("implement me")
}

func (t threadUsecase) UpdateThreadDetails(s utilities.SlugOrId, threadUpdate domain.Thread) (*domain.Thread, error) {
	panic("implement me")
}

func (t threadUsecase) GetThreadPosts(s utilities.SlugOrId, params utilities.ArrayOutParams) ([]domain.Post, error) {
	panic("implement me")
}

func (t threadUsecase) VoteThread(s utilities.SlugOrId, vote domain.Vote) (*domain.Thread, error) {
	panic("implement me")
}

func NewThreadUsecase(db *sql.DB) domain.ThreadUsecase {
	return &threadUsecase{
		DB: db,
	}
}
