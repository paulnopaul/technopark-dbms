package domain

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
