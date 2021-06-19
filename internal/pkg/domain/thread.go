package domain

import (
	"technopark-dbms/internal/pkg/utilities"
)

type Thread struct {
	ID      int32
	Title   string
	Author  string
	Forum   string
	Message string
	Votes   int32
	Slug    string
	Created string
}

type ThreadUsecase interface {
	CreatePosts(s utilities.SlugOrId, posts []Post) ([]Post, error)
	GetThreadDetails(s utilities.SlugOrId, useSlug bool) (*Thread, error)
	UpdateThreadDetails(s utilities.SlugOrId, threadUpdate Thread) (*Thread, error)
	GetThreadPosts(s utilities.SlugOrId, params utilities.ArrayOutParams) ([]Post, error)
	VoteThread(s utilities.SlugOrId, vote Vote) (*Thread, error)
}
