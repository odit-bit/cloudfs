package app

import "context"

type ctxKey string

var key ctxKey

type UserID string

func contextSetUser(ctx context.Context, userID *UserID) context.Context {
	return context.WithValue(ctx, key, userID)
}

func getUserIDFromCtx(ctx context.Context) (string, bool) {
	UserID, ok := ctx.Value(key).(*UserID)
	return string(*UserID), ok
}
