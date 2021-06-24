package forum

import (
	"errors"
	"net/http"
)

var (
	AlreadyExists   = errors.New("query creating error")
	AuthorNotExists = errors.New("query creating error")
)

func CodeFromError(err error) int {
	switch err {
	case AlreadyExists:
		return http.StatusConflict
	case AuthorNotExists:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
