package blob

import (
	"context"
	"sync"
)

var _ TokenStorer = (*memory)(nil)

type memory struct {
	mu sync.Mutex
	tm map[string]Token
}

func newObjectTokenMemory() (*memory, error) {
	return &memory{
		mu: sync.Mutex{},
		tm: map[string]Token{},
	}, nil
}

// Delete implements blob.TokenStorer.
func (m *memory) Delete(ctx context.Context, tokenKey string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.tm, tokenKey)
	return nil
}

// Get implements blob.TokenStorer.
func (m *memory) Get(ctx context.Context, tokenKey string) (*Token, bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	v, ok := m.tm[tokenKey]
	if !ok {
		return nil, false, nil
	}
	return &v, true, nil
}

// Put implements blob.TokenStorer.
func (m *memory) Put(ctx context.Context, token *Token) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.tm[token.Key()] = *token
	return nil
}
