package usecase

import "DBMSForum/internal/pkg/domain"

type usecase struct {
	forumRepo domain.ForumRepository
}

func NewForumUsecase(fr domain.ForumRepository) domain.ForumUsecase {
	return &usecase{
		forumRepo: fr,
	}
}

func (u *usecase) Create(f domain.Forum) error {
	panic("implement me")
}
