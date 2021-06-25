package forum

import (
	"errors"
	"net/http"
)

var (
	AlreadyExists   = errors.New("forum already exists")
	AuthorNotExists = errors.New("forum author does not exists")
	NotFound        = errors.New("forum not found")
)

func CodeFromError(err error) int {
	switch err {
	case AlreadyExists:
		return http.StatusConflict
	case AuthorNotExists:
		return http.StatusNotFound
	case NotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
