package blob

import "context"

type Store interface {
	Get(ctx context.Context, userID string, objname string) (*RetreiveInfo, error)
	// Put(context.Context, *putRequest) (*Result, error)
	List(ctx context.Context, conf *ListConfig) (StoreIterator, error)
}

type StoreIterator interface {
	Next() bool
	Value() *RetreiveInfo
	Error() error
}
