package domain

type Forum struct {
	Title   string
	User    string
	Slug    string
	Posts   int64
	Threads int64
}

type ForumManager interface {
	Create(f Forum) (*Forum, error)
	Details(slug string) (*Forum, error)
	CreateThread(slug string, t Thread) (*Thread, error)
	Users(slug string, limit int32, since string, desc bool) ([]User, error)
	Threads(slug string, limit int32, since string, desc bool) ([]Thread, error)
}
