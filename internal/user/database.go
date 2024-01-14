package user

import "context"

type Database interface {
	Find(ctx context.Context, name string) (*Account, error)
	Insert(ctx context.Context, account *Account) error
}
