package errs

import (
	"net/http"

	"google.golang.org/grpc/codes"
)

type ErrorCode int

const (
	ResourceInvalid ErrorCode = iota
	ResourceInUse
	ResourceNotFound
	ResourceUnauthorized
	ResourceDuplicated
)

var grpcErrorCodes = map[ErrorCode]codes.Code{
	// map code errors for gRPC codes
	ResourceInvalid:      codes.InvalidArgument,
	ResourceInUse:        codes.AlreadyExists,
	ResourceNotFound:     codes.NotFound,
	ResourceUnauthorized: codes.Unauthenticated,
	ResourceDuplicated:   codes.AlreadyExists,
}

var httpErrorCodes = map[ErrorCode]int{
	// map code errors for http codes
	ResourceInvalid:      http.StatusBadRequest,
	ResourceInUse:        http.StatusConflict,
	ResourceNotFound:     http.StatusNotFound,
	ResourceUnauthorized: http.StatusUnauthorized,
	ResourceDuplicated:   http.StatusConflict,
}
