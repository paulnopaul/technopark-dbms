package post

import "errors"

var (
	NotFoundError      = errors.New("post not found")
	InvalidParentError = errors.New("parent post was created in another thread")
)
