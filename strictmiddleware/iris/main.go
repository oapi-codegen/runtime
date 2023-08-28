package iris

import (
	"github.com/kataras/iris/v12"
)

type StrictIrisHandlerFunc func(ctx iris.Context, request interface{}) (response interface{}, err error)

type StrictIrisMiddlewareFunc func(f StrictIrisHandlerFunc, operationID string) StrictIrisHandlerFunc
