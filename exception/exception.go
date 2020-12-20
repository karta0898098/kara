package exception

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
)

var (
	// http 400
	ErrInvalidInput               = &exception{Code: 400001, Message: "One of the request inputs is not valid.", Status: http.StatusBadRequest, GRPCCode: codes.InvalidArgument}
	ErrInvalidQueryParameterValue = &exception{Code: 400002, Message: "One of the request inputs is not valid.", Status: http.StatusBadRequest, GRPCCode: codes.InvalidArgument}
	ErrInvalidHeaderValue         = &exception{Code: 400003, Message: "The value provided for one of the HTTP headers was not in the correct format.", Status: http.StatusBadRequest, GRPCCode: codes.InvalidArgument}

	// http 401
	ErrUnauthorized = &exception{Code: 401001, Message: "The request unauthorized", Status: http.StatusUnauthorized, GRPCCode: codes.PermissionDenied}

	// http 403
	ErrNotAllowed = &exception{Code: 403001, Message: "The request is understood, but it has been refused or access is not allowed.", Status: http.StatusForbidden, GRPCCode: codes.PermissionDenied}

	// http 404
	ErrPageNotFound     = &exception{Code: 404001, Message: "Page not found.", Status: http.StatusNotFound, GRPCCode: codes.NotFound}
	ErrResourceNotFound = &exception{Code: 404002, Message: "The specified resource does not exist.", Status: http.StatusNotFound}

	// 409 create resource has conflict
	ErrConflict = &exception{Code: 409001, Message: "The request conflict.", Status: http.StatusConflict, GRPCCode: codes.AlreadyExists}

	// http internal
	ErrInternal = &exception{Code: 500001, Message: "Serve occur error.", Status: http.StatusInternalServerError, GRPCCode: codes.Internal}
)

type exception struct {
	Code     int                    `json:"code"`
	Status   int                    `json:"status"`
	Message  string                 `json:"message"`
	Details  map[string]interface{} `json:"details,omitempty"`
	GRPCCode codes.Code             `json:"grpc_code"`
}

func NewException(code int, status int, message string) *exception {
	return &exception{
		Code:    code,
		Status:  status,
		Message: message,
	}
}

type exceptionView struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

func (e *exception) Error() string {
	var b strings.Builder
	_, _ = b.WriteRune('[')
	_, _ = b.WriteString(strconv.Itoa(e.Code))
	_, _ = b.WriteRune(']')
	_, _ = b.WriteRune(' ')
	_, _ = b.WriteString(e.Message)
	return b.String()
}

// As target err equal this error
func As(err error, target error) bool {
	return errors.As(err, target)
}

// Is target err equal this error
func Is(err error, target error) bool {
	return errors.Is(err, target)
}

// IsException target err equal this *exception
func IsException(target error) bool {
	_, ok := errors.Cause(target).(*exception)
	if !ok {
		return false
	}
	return true
}

// TryConvert ...
func TryConvert(target error) *exception {
	err, ok := errors.Cause(target).(*exception)
	if !ok {
		return nil
	}
	return err
}

// Is target err equal this error
func (e *exception) Is(target error) bool {
	causeErr, ok := errors.Cause(target).(*exception)
	if !ok {
		return false
	}
	return e.Code == causeErr.Code
}

// SetMessage override default message
func (e *exception) SetMessage(msg string) exception {
	e.Message = msg
	return *e
}

// SetMessagef override default message with format
func (e *exception) SetMessagef(format string, args ...interface{}) exception {
	e.Message = fmt.Sprintf(format, args...)
	return *e
}

// AddMessage add message to default message
func (e *exception) AddMessage(msg string) exception {
	var b strings.Builder
	_, _ = b.WriteString(e.Message)
	_, _ = b.WriteRune(' ')
	_, _ = b.WriteString(msg)
	e.Message = b.String()
	return *e
}

// SetDetails set detail error message
func (e *exception) SetDetails(details map[string]interface{}) exception {
	e.Details = details
	return *e
}

// AddDetails add detail error message
func (e *exception) AddDetails(key string, value interface{}) exception {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	e.Details[key] = value
	return *e
}

func (e *exception) ToView() (int, exceptionView) {
	return e.Code, exceptionView{
		Code:    e.Code,
		Message: e.Message,
		Details: e.Details,
	}
}
