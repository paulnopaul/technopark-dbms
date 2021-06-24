package errors

import (
	"encoding/json"
	"errors"
	"net/http"
)

func JSONMessage(m string) string {
	res, _ := json.Marshal(struct {
		message string
	}{m})
	return string(res)
}

func JSONErrorMessage(err error) string {
	res, _ := json.Marshal(struct {
		message error
	}{err})
	return string(res)
}

var (
	JSONEncodeErrorMessage      = JSONMessage("json encode")
	JSONDecodeErrorMessage      = JSONMessage("json decode")
	JSONURLParamsErrorMessage   = JSONMessage("url params")
	JSONQuerystringErrorMessage = JSONMessage("querystring params")
)

func CodeFromJSONMessage(message string) int {
	switch message {
	case JSONEncodeErrorMessage, JSONURLParamsErrorMessage, JSONQuerystringErrorMessage:
		return http.StatusBadRequest
	case JSONDecodeErrorMessage:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

var (
	JSONUnmarshallError   = errors.New("json unmarshall error")
	JSONEncodeError       = errors.New("json encode error")
	QuerystringParseError = errors.New("querystring parsing error")
	URLParamsError        = errors.New("url params error")
	WrongSortType         = errors.New("wrong sort type")
)

func CodeFromDeliveryError(deliveryError error) int {
	switch deliveryError {
	case URLParamsError, QuerystringParseError, JSONUnmarshallError:
		return http.StatusBadRequest
	case JSONEncodeError:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
