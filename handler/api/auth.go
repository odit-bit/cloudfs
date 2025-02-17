package api

import (
	"context"
	"net/http"

	"github.com/odit-bit/cloudfs/internal/user"
)

var (
	token_header        = "X-Token"
	token_header_expiry = "X-Token-Expiry"
)

type ctxKey int

var key ctxKey

func getTokenCtx(ctx context.Context) (*user.Token, bool) {
	sess, ok := ctx.Value(key).(*user.Token)
	return sess, ok
}

func setTokenCtx(ctx context.Context, token *user.Token) context.Context {
	return context.WithValue(ctx, key, token)
}

// non-authorized

func (v *App) Auth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	tkn, err := v.accounts.CreateToken(ctx, username, password, user.TokenOption{})
	if err != nil {
		v.serviceErr(w, r, "auth", err)
		return
	}

	w.Header().Set(token_header, tkn.Key())
	w.Header().Set(token_header_expiry, tkn.ValidUntil().Format(http.TimeFormat))
}

func (v *App) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	err := v.accounts.Register(ctx, username, password)
	if err != nil {
		v.serviceErr(w, r, "register-handler", err)
		return
	}
}

func (v *App) NeedToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get(token_header)
		if token == "" {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		tkn, err := v.accounts.TokenAuth(r.Context(), token)
		if err != nil {
			v.serviceErr(w, r, "ApiMiddleware", err)
			return
		}

		ctx := setTokenCtx(r.Context(), tkn)
		rr := r.WithContext(ctx)
		next.ServeHTTP(w, rr)

		// if v.session.Status(ctx) == scs.Modified {
		// 	token, exp, err := v.session.Commit(ctx)
		// 	if err != nil {
		// 		v.serviceErr(w, r, "ApiMidlleware", err)
		// 		return
		// 	}
		// 	w.Header().Set(token_header, token)
		// 	w.Header().Set(token_header_expiry, exp.Format(http.TimeFormat))
		// }

	})
}

func (v *App) CreateShareToken(w http.ResponseWriter, r *http.Request) {
	tkn, _ := getTokenCtx(r.Context())
	filename := r.URL.Query().Get("filename")
	if filename == "" {
		http.Error(w, "filename is empty", http.StatusBadRequest)
	}
	shareToken, err := v.objects.CreateShareToken(r.Context(), tkn.UserID(), filename, user.Default_Token_Expire)
	if err != nil {
		v.serviceErr(w, r, "createFileShareToken", err)
		return
	}

	w.Write([]byte(shareToken.Key()))
}
