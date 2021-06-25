package domain

import (
	"technopark-dbms/internal/pkg/utilities"
)

type Forum struct {
	Title   string `json:"title"`
	User    string `json:"user"`
	Slug    string `json:"slug"`
	Posts   int64  `json:"posts,omitempty"`
	Threads int64  `json:"threads,omitempty"`
}

type ForumUsecase interface {
	CreateForum(f Forum) (*Forum, error)
	Exists(slug string) (bool, error)
	Details(slug string) (*Forum, error)
	CreateThread(forumSlug string, t Thread) (*Thread, error)
	Users(slug string, params utilities.ArrayOutParams) ([]User, error)
	Threads(forumSlug string, params utilities.ArrayOutParams) ([]Thread, error)
}

type Post struct {
	ID       int64  `json:"id"`
	Parent   int64  `json:"parent,omitempty"`
	Author   string `json:"author,omitempty"`
	Message  string `json:"message,omitempty"`
	IsEdited bool   `json:"isEdited,omitempty"`
	Forum    string `json:"forum,omitempty"`
	Thread   int32  `json:"thread,omitempty"`
	Created  string `json:"created,omitempty"`
}

type PostUsecase interface {
	GetById(id int64) (*Post, error)
	GetDetails(id int64, relatedUser bool, relatedForum bool, relatedThread bool) (*Post, *Forum, *Thread, *User, error)
	UpdateDetails(id int64, postUpdate Post) (*Post, error)
}

type Service struct {
	User   int32 `json:"user"`
	Forum  int32 `json:"forum"`
	Thread int32 `json:"thread"`
	Post   int64 `json:"post"`
}

type ServiceUsecase interface {
	Clear() error
	Status() (*Service, error)
}

type Thread struct {
	ID      int32  `json:"id"`
	Title   string `json:"title"`
	Author  string `json:"author"`
	Forum   string `json:"forum"`
	Message string `json:"message"`
	Votes   int32  `json:"votes,omitempty"`
	Slug    string `json:"slug,omitempty"`
	Created string `json:"created,omitempty"`
}

type ThreadUsecase interface {
	CreatePosts(s utilities.SlugOrId, posts []Post) ([]Post, error)
	GetThreadDetails(s utilities.SlugOrId) (*Thread, error)
	UpdateThreadDetails(s utilities.SlugOrId, threadUpdate Thread) (*Thread, error)
	GetThreadPosts(s utilities.SlugOrId, params utilities.ArrayOutParams) ([]Post, error)
	VoteThread(s utilities.SlugOrId, vote Vote) (*Thread, error)
}

type User struct {
	Nickname string `json:"nickname,omitempty"`
	Fullname string `json:"fullname,omitempty"`
	About    string `json:"about,omitempty"`
	Email    string `json:"email,omitempty"`
}

type UserUsecase interface {
	CreateUser(nickname string, createData User) (*User, error, []User)
	GetProfile(nickname string) (*User, error)
	GetProfiles(nickname, email string) ([]User, error)
	UpdateProfile(nickname string, profileUpdate User) (*User, error)
	Exists(nickname string, email string) (bool, error)
}

type Vote struct {
	Nickname string `json:"nickname"`
	Voice    int32  `json:"voice"`
}
