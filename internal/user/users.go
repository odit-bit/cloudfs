package user

import (
	"context"
	"errors"
	"fmt"
)

type AccountStorer interface {
	FindUsername(ctx context.Context, username string) (*Account, error)
	FindID(ctx context.Context, id string) (*Account, bool, error)
	Insert(ctx context.Context, acc *Account) error
}

type TokenStorer interface {
	GetToken(ctx context.Context, tkn string) (*Token, error)
	GetTokenUserID(ctx context.Context, id string) (*Token, bool)
	PutToken(ctx context.Context, token *Token) error
	Delete(ctx context.Context, tkn string) error
}

type Users struct {
	accounts AccountStorer
	tokens   TokenStorer
}

func NewWithMemory() *Users {
	db, _ := newInMemory()
	return &Users{
		accounts: db,
		tokens:   db,
	}
}

func New(ctx context.Context, accounts AccountStorer, tokens TokenStorer) (*Users, error) {
	st := Users{accounts: accounts, tokens: tokens}
	return &st, nil
}

func (s *Users) Register(ctx context.Context, username, password string) (*Account, error) {
	_, err := s.accounts.FindUsername(ctx, username)
	if err == nil {
		return nil, errors.Join(ErrAccountExist)
	}
	acc := CreateAccount(username, password)
	if err := s.accounts.Insert(ctx, acc); err != nil {
		return nil, err
	}
	return acc, nil
}

type BasicAuthResponse struct {
	ID string
	*Token
}

func (s *Users) BasicAuth(ctx context.Context, username, password string) (*BasicAuthResponse, error) {
	acc, err := s.accounts.FindUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	if ok := acc.CheckPassword(password); !ok {
		return nil, err
	}

	tkn, ok := s.tokens.GetTokenUserID(ctx, acc.ID.String())
	if !ok {
		tkn = NewToken(acc.ID.String(), Default_Token_Expire)
	}

	if err := s.tokens.PutToken(ctx, tkn); err != nil {
		return nil, err
	}

	return &BasicAuthResponse{ID: acc.ID.String(), Token: tkn}, nil
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

func (s *Users) GetAccount(ctx context.Context, id string) (*Account, bool, error) {
	return s.accounts.FindID(ctx, id)
}
