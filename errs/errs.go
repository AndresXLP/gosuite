package errs

import (
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AppError struct {
	error
	Code    ErrorCode
	Message string
}

func (a *AppError) getHttpCode() int {
	return httpErrorCodes[a.Code]
}

func (a *AppError) getGRPCCode() codes.Code {
	return grpcErrorCodes[a.Code]
}

func (a *AppError) Error() string {
	return a.Message
}

// NewAppError creates a new AppError instance.
func NewAppError(code ErrorCode, msg string) *AppError {
	he := &AppError{
		Code:    code,
		Message: msg,
	}

	return he
}

// NewEchoHttpError this method return a new instance HttpError
func (a *AppError) NewEchoHttpError() error {
	return echo.NewHTTPError(a.getHttpCode(), a.Error())
}

func (a *AppError) NewGRPCError() error {
	return status.Error(a.getGRPCCode(), a.Error())
}
