package errors

import (
	"net/http"

	"google.golang.org/grpc/codes"
)

var (
	// http 400
	ErrInvalidInput       = &Exception{Code: 400001, Message: "One of the request inputs is not valid.", Status: http.StatusBadRequest, GRPCCode: codes.InvalidArgument}
	ErrInvalidHeaderValue = &Exception{Code: 400003, Message: "The value provided for one of the HTTP headers was not in the correct format.", Status: http.StatusBadRequest, GRPCCode: codes.InvalidArgument}

	// http 401
	ErrUnauthorized = &Exception{Code: 401001, Message: "The request unauthorized", Status: http.StatusUnauthorized, GRPCCode: codes.PermissionDenied}

	// http 403
	ErrNotAllowed = &Exception{Code: 403001, Message: "The request is understood, but it has been refused or access is not allowed.", Status: http.StatusForbidden, GRPCCode: codes.PermissionDenied}

	// http 404
	ErrPageNotFound     = &Exception{Code: 404001, Message: "Page not found.", Status: http.StatusNotFound, GRPCCode: codes.NotFound}
	ErrResourceNotFound = &Exception{Code: 404002, Message: "The specified resource does not exist.", Status: http.StatusNotFound}

	// 409 create resource has conflict
	ErrConflict = &Exception{Code: 409001, Message: "The request conflict.", Status: http.StatusConflict, GRPCCode: codes.AlreadyExists}

	// http 429 too many request
	ErrTooManyRequests = &Exception{Code: 429001, Message: "Too Many Requests", Status: http.StatusTooManyRequests, GRPCCode: codes.PermissionDenied}

	// http internal
	ErrInternal = &Exception{Code: 500001, Message: "Serve occur error.", Status: http.StatusInternalServerError, GRPCCode: codes.Internal}
)
