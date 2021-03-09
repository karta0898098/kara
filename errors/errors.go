package errors

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
)

type DetailData map[string]interface{}

type Exception struct {
	Code     int        `json:"code"`
	Status   int        `json:"status"`
	Message  string     `json:"message"`
	Details  DetailData `json:"details,omitempty"`
	GRPCCode codes.Code `json:"grpc_code"`
}

// New server internal error with message
func New(message string) *Exception {
	return &Exception{
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
func TryConvert(target error) *Exception {
	err, ok := errors.Cause(target).(*Exception)
	if !ok {
		return nil
	}
	return err
}

// Is Check input is same
func Is(err error, target error) bool {
	return errors.Is(err, target)
}

func (e *Exception) Error() string {
	var b strings.Builder
	_, _ = b.WriteRune('[')
	_, _ = b.WriteString(strconv.Itoa(e.Code))
	_, _ = b.WriteRune(']')
	_, _ = b.WriteRune(' ')
	_, _ = b.WriteString(e.Message)
	return b.String()
}

// Is target err equal this error
func (e *Exception) Is(target error) bool {
	causeErr, ok := errors.Cause(target).(*Exception)
	if !ok {
		return false
	}
	return e.Code == causeErr.Code
}

// WithDetails set detail error message
func (e *Exception) WithDetails(details DetailData) *Exception {
	newErr := *e
	newErr.Details = details
	return &newErr
}

func (e *Exception) BuildWithError(err error) error {
	return errors.Wrap(e, err.Error())
}

func (e *Exception) Build(format string, args ...interface{}) error {
	return errors.Wrapf(e, format, args...)
}

// ToViewModel to restful view
func (e *Exception) ToViewModel() (int, *exceptionView) {
	return e.Status, &exceptionView{
		Code:    e.Code,
		Message: e.Message,
		Details: e.Details,
	}
}
