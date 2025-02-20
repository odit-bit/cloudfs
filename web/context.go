package web

import "context"

const (
	sessUserToken = "userToken"
)

type ctxKey string

var key ctxKey

type userToken string

func setUserTokenCtx(ctx context.Context, userToken *userToken) context.Context {
	return context.WithValue(ctx, key, userToken)
}

func getUserTokenFromCtx(ctx context.Context) (string, bool) {
	UserID, ok := ctx.Value(key).(*userToken)
	return string(*UserID), ok
}
