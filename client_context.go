package runtime

import "context"

type operationIdKey int

const key = operationIdKey(0)

func WithOperationId(ctx context.Context, operationId string) context.Context {
	return context.WithValue(ctx, key, operationId)
}

func GetOperationId(ctx context.Context) string {
	if operationId, ok := ctx.Value(key).(string); ok {
		return operationId
	}
	return ""
}
