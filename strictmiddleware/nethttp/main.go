package nethttp

import (
	"context"
	"net/http"
)

type StrictHttpHandlerFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (response interface{}, err error)

type StrictHttpMiddlewareFunc func(f StrictHttpHandlerFunc, operationID string) StrictHttpHandlerFunc
