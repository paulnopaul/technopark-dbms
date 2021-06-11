package errors

import (
	"errors"
	"fmt"
	"net/http"
)

func JSONMessage(m string) string {
	return fmt.Sprintf(`{"message":"%s"}`, m)
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
