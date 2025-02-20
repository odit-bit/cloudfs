package user

import (
	"context"
	"sync"
)

var _ AccountStorer = (*memory)(nil)
var _ TokenStorer = (*memory)(nil)

type memory struct {
	mu      sync.Mutex
	m       map[string]Account
	mTkn    map[string]Token
	indexID map[string]string // map[userID]tokenKey
}

// only for testing
func newInMemory() (*memory, error) {
	return &memory{
		mu:      sync.Mutex{},
		m:       map[string]Account{},
		mTkn:    map[string]Token{},
		indexID: map[string]string{},
	}, nil
}

// Delete implements TokenStorer
func (i *memory) Delete(ctx context.Context, tkn string) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	delete(i.mTkn, tkn)
	return nil
}

// GetToken implements TokenStorer.
func (i *memory) GetToken(ctx context.Context, tkn string) (*Token, error) {
	i.mu.Lock()
	defer i.mu.Unlock()
	v, ok := i.mTkn[tkn]
	if !ok {
		return nil, ErrTokenNotExist
	}
	return &v, nil
}

// PutToken implements TokenStorer.
func (i *memory) PutToken(ctx context.Context, token *Token) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	oldToken, ok := i.indexID[token.UserID()]
	if ok {
		delete(i.mTkn, oldToken)
	}

	i.indexID[token.UserID()] = token.Key()
	i.mTkn[token.Key()] = *token

	return nil
}

func (i *memory) GetTokenUserID(ctx context.Context, id string) (*Token, bool) {
	i.mu.Lock()
	defer i.mu.Unlock()
	key, ok := i.indexID[id]
	if !ok {
		return nil, false
	}
	tkn, ok := i.mTkn[key]
	return &tkn, ok
}

// Find implements service.AccountStore.
func (i *memory) FindUsername(ctx context.Context, username string) (*Account, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	u, ok := i.m[username]
	if !ok {
		return nil, ErrAccountNotExist
	}

	result := Account{
		ID:           u.ID,
		Name:         u.Name,
		HashPassword: u.HashPassword,
		Quota:        u.Quota,
	}
	return &result, nil

}

// Insert implements service.AccountStore.
func (i *memory) Insert(ctx context.Context, acc *Account) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	_, ok := i.m[acc.Name]
	if ok {
		return ErrAccountExist
	}

	i.m[acc.Name] = *acc
	return nil
}
