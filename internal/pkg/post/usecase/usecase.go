package usecase

import (
	"DBMSForum/internal/pkg/domain"
	"database/sql"
)

type postUsecase struct {
	DB *sql.DB
}

func NewPostUsecase(db *sql.DB) domain.PostUsecase {
	return &postUsecase{
		DB: db,
	}
}

func (p *postUsecase) GetDetails(id int64, relatedUser bool, relatedForum bool, relatedThread bool) (*domain.Post, error) {
	panic("implement me")
}

func (p *postUsecase) UpdateDetails(id int64, postUpdate domain.Post) (*domain.Post, error) {
	panic("implement me")
}
