package domain

import "DBMSForum/internal/pkg/utilities"

type Forum struct {
	Title   string
	User    string
	Slug    string
	Posts   int64
	Threads int64
}

type ForumUsecase interface {
	Create(f Forum) (*Forum, error)
	Details(slug string) (*Forum, error)
	CreateThread(slug string, t Thread) (*Thread, error)
	Users(slug string, params utilities.ArrayOutParams) ([]User, error)
	Threads(slug string, params utilities.ArrayOutParams) ([]Thread, error)
}
