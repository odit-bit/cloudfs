package repo

import (
	"context"
	"fmt"
	"sync"

	"github.com/odit-bit/cloudfs/internal/user"
	"github.com/odit-bit/cloudfs/service"
)

var _ service.AccountStore = (*memory)(nil)

type memory struct {
	mu sync.Mutex
	m  map[string]user.Account
}

func NewInMemory() (*memory, error) {
	return &memory{
		mu: sync.Mutex{},
		m:  map[string]user.Account{},
	}, nil
}

// Find implements service.AccountStore.
func (i *memory) Find(ctx context.Context, username string) (*user.Account, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	u, ok := i.m[username]
	if !ok {
		return nil, fmt.Errorf("not found")
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
		return fmt.Errorf("username exist")
	}

	i.m[acc.Name] = *acc
	return nil
}
