package errors

import (
	"errors"
	"fmt"
	"net/http"
	"technopark-dbms/internal/pkg/domain"
)

func JSONMessage(m string) domain.JSONMessageType {
	return domain.JSONMessageType{Message: m}
}

func JSONErrorMessage(err error) domain.JSONMessageType {
	return domain.JSONMessageType{Message: fmt.Sprint(err)}
}

var (
	JSONEncodeErrorMessage      = JSONMessage("json encode")
	JSONDecodeErrorMessage      = JSONMessage("json decode")
	JSONURLParamsErrorMessage   = JSONMessage("url params")
	JSONQuerystringErrorMessage = JSONMessage("querystring params")
)

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
