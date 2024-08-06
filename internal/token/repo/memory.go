package repo

import (
	"context"
	"sync"

	"github.com/odit-bit/cloudfs/internal/token"
	"github.com/odit-bit/cloudfs/service"
)

var _ service.TokenStore = (*memory)(nil)

type memory struct {
	l      sync.Mutex
	tokens map[string]token.ShareToken
}

func NewInMemory() (*memory, error) {
	return &memory{
		l:      sync.Mutex{},
		tokens: map[string]token.ShareToken{},
	}, nil
}

// Get implements service.TokenStore.
func (m *memory) Get(ctx context.Context, tokenString string) (*token.ShareToken, bool, error) {
	m.l.Lock()
	defer m.l.Unlock()
	v, ok := m.tokens[tokenString]
	if !ok {
		return nil, false, nil
	}
	return &v, true, nil
}

// Put implements service.TokenStore.
func (m *memory) Put(ctx context.Context, token *token.ShareToken) error {
	m.l.Lock()
	defer m.l.Unlock()
	m.tokens[token.Key()] = *token
	return nil
}
