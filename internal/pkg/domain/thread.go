package domain

import "strconv"

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

type ArrayOutParams struct {
	Limit    int32
	HasLimit bool
	Since    string
	HasSince bool
	Desc     bool
	Sort     string
	HasSort  bool
}

type SlugOrId struct {
	IsSlug bool
	Slug   string
	ID     int32
}

func NewSlugOrId(idString string) SlugOrId {
	parsedID, err := strconv.ParseInt(idString, 10, 32)
	if err != nil {
		return SlugOrId{
			IsSlug: true,
			Slug:   idString,
		}
	}
	return SlugOrId{
		IsSlug: false,
		ID:     int32(parsedID),
	}
}

type ThreadUsecase interface {
	CreatePosts(s SlugOrId, posts []Post) ([]Post, error)
	GetThreadDetails(s SlugOrId, useSlug bool) (*Thread, error)
	UpdateThreadDetails(s SlugOrId, threadUpdate Thread) (*Thread, error)
	GetThreadPosts(s SlugOrId, params ArrayOutParams) ([]Post, error)
	VoteThread(s SlugOrId, vote Vote) (*Thread, error)
}
