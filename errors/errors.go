package errors

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
)

type exception struct {
	Code     int                    `json:"code"`
	Status   int                    `json:"status"`
	Message  string                 `json:"message"`
	Details  map[string]interface{} `json:"details,omitempty"`
	GRPCCode codes.Code             `json:"grpc_code"`
}

// New server internal error with message
func New(message string) *exception {
	return &exception{
		Code:     500001,
		Status:   http.StatusInternalServerError,
		Message:  message,
		GRPCCode: codes.Internal,
	}
}

type exceptionView struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// TryConvert ...
func TryConvert(target error) *exception {
	err, ok := target.(*exception)
	if !ok {
		return nil
	}
	return err
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

// Is target err equal this error
func (e *exception) Is(target error) bool {
	causeErr, ok := errors.Cause(target).(*exception)
	if !ok {
		return false
	}
	return e.Code == causeErr.Code
}

// SetMessage override default message
func (e *exception) SetMessage(msg string) *exception {
	e.Message = msg
	return e
}

// SetMessagef override default message with format
func (e *exception) SetMessagef(format string, args ...interface{}) *exception {
	e.Message = fmt.Sprintf(format, args...)
	return e
}

// AddMessage add message to default message
func (e *exception) AddMessage(msg string) *exception {
	var b strings.Builder
	_, _ = b.WriteString(e.Message)
	_, _ = b.WriteRune(',')
	_, _ = b.WriteString(msg)
	e.Message = b.String()
	return e
}

// SetDetails set detail error message
func (e *exception) SetDetails(details map[string]interface{}) *exception {
	e.Details = details
	return e
}

// AddDetails add detail error message
func (e *exception) AddDetails(key string, value interface{}) *exception {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	e.Details[key] = value
	return e
}

// ToRestfulView to restful view
func (e *exception) ToRestfulView() (int, *exceptionView) {
	return e.Code, &exceptionView{
		Code:    e.Code,
		Message: e.Message,
		Details: e.Details,
	}
}
