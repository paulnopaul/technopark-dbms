package db

import "DBMSForum/internal/pkg/domain"

type repo struct {
}

func NewForumRepository() domain.ForumRepository {
	return &repo{}
}

func (r *repo) Create(f domain.Forum) error {
	panic("implement me")
}
