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

func NewAppError(code ErrorCode, msg string) *AppError {
	return &AppError{
		Code:    code,
		Message: msg,
	}
}

// NewEchoHttpError this method return echo.NewHttpError
func (a *AppError) NewEchoHttpError() error {
	return echo.NewHTTPError(a.getHttpCode(), a.Error())
}

func (a *AppError) NewGRPCError() error {
	return status.Error(a.getGRPCCode(), a.Error())
}
