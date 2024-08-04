package service

import (
	"context"

	"github.com/odit-bit/cloudfs/internal/user"
)

type AuthParam struct {
	Username string
	Password string
}

type AuthResult struct {
	ID string
}

func (app *Cloudfs) Auth(ctx context.Context, param *AuthParam) (*AuthResult, error) {

	acc, err := app.accounts.Find(ctx, param.Username)
	if err != nil {
		return nil, err
	}
	if ok := acc.CheckPassword(param.Password); !ok {
		return nil, ErrInvalidCredentials
	}

	return &AuthResult{
		ID: acc.ID.String(),
	}, err
}

type RegisterParam struct {
	Username string
	Password string
}

func (app *Cloudfs) Register(ctx context.Context, param *RegisterParam) error {
	if _, err := app.accounts.Find(ctx, param.Username); err != nil {
		err = nil
		acc := user.CreateAccount(param.Username, param.Password)
		return app.accounts.Insert(ctx, acc)
	}

	return ErrAccountExist
}
