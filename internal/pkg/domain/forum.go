package domain

type Forum struct {
	Title   string
	User    string
	Slug    string
	Posts   int64
	Threads int64
}

type ForumUsecase interface {
	Create(f Forum) error
}

type ForumRepository interface {
	Create(f Forum) error
}
