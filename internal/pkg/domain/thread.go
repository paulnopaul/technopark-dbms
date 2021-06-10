package domain

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
	CreateThread()
	GetThreadDetails()
	UpdateThreadDetails()
	GetThreadPosts()
	VoteThread()
}
