package factory

import (
	"fmt"
	"net/http"
)

const WrapErrorMessage = "echosample error"

type Error struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Details string `json:"details,omitempty"`
	err     error
	status  int
}

func (t ErrorTemplate) New(err error, v ...interface{}) *Error {
	e := Error{
		Code:    t.Code,
		Message: fmt.Sprintf(t.Message, v...),
		err:     err,
	}
	if err != nil {
		e.Details = fmt.Sprintf("%s: %s", WrapErrorMessage, err.Error())
	}
	return &e
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	return e.Details
}

func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.err
}

func (e *Error) Status() int {
	if e == nil || e.status == 0 {
		return http.StatusInternalServerError
	}
	return e.status
}

type ErrorTemplate Error

var (
	// System Error
	ErrorSystem             = ErrorTemplate{Code: 10001, Message: "System Error"}
	ErrorServiceUnavailable = ErrorTemplate{Code: 10002, Message: "Service unavailable"}
	ErrorRemoteService      = ErrorTemplate{Code: 10003, Message: "Remote service error"}
	ErrorIPLimit            = ErrorTemplate{Code: 10004, Message: "IP limit"}
	ErrorPermissionDenied   = ErrorTemplate{Code: 10005, Message: "Permission denied", status: http.StatusForbidden}
	ErrorIllegalRequest     = ErrorTemplate{Code: 10006, Message: "Illegal request", status: http.StatusBadRequest}
	ErrorHTTPMethod         = ErrorTemplate{Code: 10007, Message: "HTTP method is not suported for this request", status: http.StatusMethodNotAllowed}
	ErrorParameter          = ErrorTemplate{Code: 10008, Message: "Parameter error", status: http.StatusBadRequest}
	ErrorMissParameter      = ErrorTemplate{Code: 10009, Message: "Miss required parameter", status: http.StatusBadRequest}
	ErrorDB                 = ErrorTemplate{Code: 10010, Message: "DB error, please contact the administator"}
	ErrorTokenInvaild       = ErrorTemplate{Code: 10011, Message: "Token invaild", status: http.StatusUnauthorized}
	ErrorMissToken          = ErrorTemplate{Code: 10012, Message: "Miss token", status: http.StatusUnauthorized}
	ErrorVersion            = ErrorTemplate{Code: 10013, Message: "API version %s invalid"}
	ErrorNotFound           = ErrorTemplate{Code: 10014, Message: "Resource not found", status: http.StatusNotFound}
	// Business Error
	ErrorDiscountNotExists = ErrorTemplate{Code: 20001, Message: "Discount %d does not exists"}
)
