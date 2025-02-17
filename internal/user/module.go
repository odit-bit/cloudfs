package user

import (
	"context"
	"errors"
	"fmt"
)

type AccountStorer interface {
	FindUsername(ctx context.Context, username string) (*Account, error)
	Insert(ctx context.Context, acc *Account) error
}

type TokenStorer interface {
	GetToken(ctx context.Context, tkn string) (*Token, error)
	PutToken(ctx context.Context, token *Token) error
	Delete(ctx context.Context, tkn string) error
}

type Users struct {
	accounts AccountStorer
	tokens   TokenStorer
}

func NewStore(ctx context.Context, accounts AccountStorer, tokens TokenStorer) (*Users, error) {
	st := Users{accounts: accounts, tokens: tokens}
	return &st, nil
}

func (s *Users) Register(ctx context.Context, username, password string) error {
	acc := CreateAccount(username, password)
	return s.accounts.Insert(ctx, acc)
}

func (s *Users) BasicAuth(ctx context.Context, username, password string) (*Account, error) {
	acc, err := s.accounts.FindUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	if ok := acc.CheckPassword(password); !ok {
		return nil, err
	}

	return acc, nil
}

func (s *Users) TokenAuth(ctx context.Context, tkn string) (*Token, error) {
	v, err := s.tokens.GetToken(ctx, tkn)
	if err != nil {
		return nil, err
	}

	if !v.IsNotExpire() {
		err2 := s.tokens.Delete(ctx, tkn)
		return nil, errors.Join(err2, ErrTokenExpired)
	}

	// not expired, refresh time expire
	v.RefreshExpire(Default_Token_Expire)
	if err := s.tokens.PutToken(ctx, v); err != nil {
		return nil, err
	}

	return v, nil

}

func (s *Users) CreateToken(ctx context.Context, username, password string, opts TokenOption) (*Token, error) {
	acc, err := s.accounts.FindUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	if ok := acc.CheckPassword(password); !ok {
		return nil, fmt.Errorf("invalid credentials")
	}

	// create token
	tkn := NewToken(acc.ID.String(), opts.Expire)

	if err := s.tokens.PutToken(ctx, tkn); err != nil {
		return nil, err
	}

	return tkn, nil
}
