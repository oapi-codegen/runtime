package v5

import "github.com/labstack/echo/v5"

type StrictEchoHandlerFunc func(ctx *echo.Context, request interface{}) (response interface{}, err error)

type StrictEchoMiddlewareFunc func(f StrictEchoHandlerFunc, operationID string) StrictEchoHandlerFunc
