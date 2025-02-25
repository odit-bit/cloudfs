package web

import (
	"context"
	"encoding/gob"
	"time"
)

// const (
// 	sessUserToken = "userToken"
// )

// type ctxKey string

// var key ctxKey

// type userToken string

// func setUserTokenCtx(ctx context.Context, userToken *userToken) context.Context {
// 	return context.WithValue(ctx, key, userToken)
// }

// func getUserTokenFromCtx(ctx context.Context) (string, bool) {
// 	UserID, ok := ctx.Value(key).(*userToken)
// 	return string(*UserID), ok
// }

//// session account

const (
	Session_Account = "session_account"
)

func init() {
	gob.Register(&account{})

}

type accKey int

var accKeyCtx accKey

type Filename string

type account struct {
	UserID        string
	Token         string
	SharedObjects map[Filename]ShareToken
}

type ShareToken struct {
	Key        string
	ValidUntil time.Time
}

// type SharedObjectToken map[filename]tokenString

func getAccountFromCtx(ctx context.Context) (*account, bool) {
	acc, ok := ctx.Value(accKeyCtx).(*account)
	return acc, ok
}

func setAccountCtx(ctx context.Context, acc *account) context.Context {
	return context.WithValue(ctx, accKeyCtx, acc)
}
