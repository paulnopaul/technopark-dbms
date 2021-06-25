package thread

import "errors"

var (
	AlreadyExists   = errors.New("thread already exists")
	NotFound        = errors.New("thread not found")
	AuthorNotExists = errors.New("thread author does not exist")
)
