package errors

import (
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"strings"
)

var (
	ErrInvalidInput     = &AppError{Code: 400001, Message: "One of the request inputs is not valid.", Status: http.StatusBadRequest}
	ErrUnauthorized     = &AppError{Code: 401001, Message: http.StatusText(http.StatusUnauthorized), Status: http.StatusUnauthorized}
	ErrResourceNotFound = &AppError{Code: 404001, Message: "The specified resource does not exist.", Status: http.StatusNotFound}
)

type AppError struct {
	Code    int    `json:"code"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (e *AppError) Error() string {
	var b strings.Builder
	_, _ = b.WriteRune('[')
	_, _ = b.WriteString(strconv.Itoa(e.Code))
	_, _ = b.WriteRune(']')
	_, _ = b.WriteRune(' ')
	_, _ = b.WriteString(e.Message)
	return b.String()
}

func (e *AppError) Is(target error) bool {

	causeErr, ok := errors.Cause(target).(*AppError)
	if !ok {
		return false
	}
	return e.Code == causeErr.Code
}
