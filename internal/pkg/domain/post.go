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
