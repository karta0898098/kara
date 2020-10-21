package exception

import (
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"net/http"
	"strconv"
	"strings"
)

var (
	ErrInvalidInput               = &Exception{Code: 400001, Message: "One of the request inputs is not valid.", Status: http.StatusBadRequest, GRPCCode: codes.InvalidArgument}
	ErrInvalidQueryParameterValue = &Exception{Code: 400002, Message: "One of the request inputs is not valid.", Status: http.StatusBadRequest}
	ErrInvalidHeaderValue         = &Exception{Code: 400003, Message: "The value provided for one of the HTTP headers was not in the correct format.", Status: http.StatusBadRequest}

	ErrUnauthorized = &Exception{Code: 401001, Message: http.StatusText(http.StatusUnauthorized), Status: http.StatusUnauthorized}

	ErrNotAllowed   = &Exception{Code: 403001, Message: "The request is understood, but it has been refused or access is not allowed.", Status: http.StatusForbidden}
	ErrPageNotFound = &Exception{Code: 404002, Message: "Page Not Found.", Status: http.StatusNotFound}

	ErrResourceNotFound = &Exception{Code: 404001, Message: "The specified resource does not exist.", Status: http.StatusNotFound}
	ErrServerInternal   = &Exception{Code: 500001, Message: http.StatusText(http.StatusInternalServerError), Status: http.StatusInternalServerError}
)

type Exception struct {
	Code     int                    `json:"code"`
	Status   int                    `json:"status"`
	Message  string                 `json:"message"`
	Details  map[string]interface{} `json:"details"`
	GRPCCode codes.Code             `json:"grpc_code"`
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

func (e *Exception) Is(target error) bool {

	causeErr, ok := errors.Cause(target).(*Exception)
	if !ok {
		return false
	}
	return e.Code == causeErr.Code
}
