package domain

import "technopark-dbms/internal/pkg/utilities"

type Forum struct {
	Title   string
	User    string
	Slug    string
	Posts   int64
	Threads int64
}

type ForumUsecase interface {
	CreateForum(f Forum) (*Forum, error)
	Details(slug string) (*Forum, error)
	CreateThread(slug string, t Thread) (*Thread, error)
	Users(slug string, params utilities.ArrayOutParams) ([]User, error)
	Threads(slug string, params utilities.ArrayOutParams) ([]Thread, error)
}

type Post struct {
	ID       int64
	Parent   int64
	Author   string
	Message  string
	IsEdited bool `json:"idEdited"`
	Forum    string
	Thread   string
	Created  string
}

type PostUsecase interface {
	GetDetails(id int64, relatedUser bool, relatedForum bool, relatedThread bool) (*Post, error)
	UpdateDetails(id int64, postUpdate Post) (*Post, error)
}

type Service struct {
	User   int32
	Forum  int32
	Thread int32
	Post   int64
}

type ServiceUsecase interface {
	Clear() error
	Status() (*Service, error)
}

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

type User struct {
	Nickname string
	Fullname string
	About    string
	Email    string
}

type UserUsecase interface {
	CreateUser(nickname string, createData User) (*User, error)
	GetProfile(nickname string) (*User, error)
	UpdateProfile(nickname string, profileUpdate User) (*User, error)
	Exists(nickname string, email string) (bool, error)
}

type Vote struct {
	Nickname string
	Voice    int32
}
