package user

import (
	"errors"
	"net/http"
)

var (
	AlreadyExistsError = errors.New("user already exists")
	NotExistsError     = errors.New("user does not exist")
	UpdateConflict     = errors.New("user update conflicts with another users")
)

func CodeFromError(err error) int {
	switch err {
	case AlreadyExistsError, UpdateConflict:
		return http.StatusConflict
	case NotExistsError:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
