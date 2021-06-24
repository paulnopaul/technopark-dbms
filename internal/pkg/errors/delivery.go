package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type JSONMessageType struct {
	Message string `json:"message"`
}

func JSONMessage(m string) []byte {
	res, _ := json.Marshal(JSONMessageType{m})
	return res
}

func JSONErrorMessage(err error) []byte {
	res, _ := json.Marshal(JSONMessageType{fmt.Sprint(err)})
	return res
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
