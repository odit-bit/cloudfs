package repo

import (
	"context"
	"sync"

	"github.com/odit-bit/cloudfs/internal/user"
)

var _ user.AccountStorer = (*memory)(nil)
var _ user.TokenStorer = (*memory)(nil)

type memory struct {
	mu   sync.Mutex
	m    map[string]user.Account
	mTkn map[string]user.Token
}

func NewInMemory() (*memory, error) {
	return &memory{
		mu:   sync.Mutex{},
		m:    map[string]user.Account{},
		mTkn: map[string]user.Token{},
	}, nil
}

// Delete implements user.TokenStorer
func (i *memory) Delete(ctx context.Context, tkn string) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	delete(i.mTkn, tkn)
	return nil
}

// GetToken implements user.TokenStorer.
func (i *memory) GetToken(ctx context.Context, tkn string) (*user.Token, error) {
	i.mu.Lock()
	defer i.mu.Unlock()
	v, ok := i.mTkn[tkn]
	if !ok {
		return nil, user.ErrTokenNotExist
	}
	return &v, nil
}

// PutToken implements user.TokenStorer.
func (i *memory) PutToken(ctx context.Context, token *user.Token) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.mTkn[token.Key()] = *token
	return nil
}

// Find implements service.AccountStore.
func (i *memory) FindUsername(ctx context.Context, username string) (*user.Account, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	u, ok := i.m[username]
	if !ok {
		return nil, user.ErrAccountNotExist
	}

	result := user.Account{
		ID:           u.ID,
		Name:         u.Name,
		HashPassword: u.HashPassword,
		Quota:        u.Quota,
	}
	return &result, nil

}

// Insert implements service.AccountStore.
func (i *memory) Insert(ctx context.Context, acc *user.Account) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	_, ok := i.m[acc.Name]
	if ok {
		return user.ErrAccountExist
	}

	i.m[acc.Name] = *acc
	return nil
}
