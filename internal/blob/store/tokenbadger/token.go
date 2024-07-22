package tokenbadger

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"log"
	"strings"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/odit-bit/cloudfs/service"
)

var _ service.TokenStore = (*Manager)(nil)

type Manager struct {
	db *badger.DB
}

// Query implements service.TokenStore.
func (m *Manager) Query(ctx context.Context, txn func(txn service.TokenTxn) error) error {
	panic("unimplemented")
}

func (m *Manager) Close() error {
	return m.db.Close()
}

func connectBadger(path string) *badger.DB {
	opt := badger.DefaultOptions(path)
	if path == "" {
		opt = opt.WithInMemory(true)
	}
	bdb, err := badger.Open(opt)
	if err != nil {
		log.Fatal(err)
	}

	return bdb

}

func New(path string) (*Manager, error) {
	return &Manager{
		db: connectBadger(path),
	}, nil
}

// Generate implements service.TokenStore.
func (m *Manager) Generate(ctx context.Context, bucket string, filename string, dur time.Duration) (string, error) {

	//make token from this enc
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	tokenString := hex.EncodeToString(b)
	value := base64.URLEncoding.EncodeToString([]byte(strings.Join([]string{bucket, filename}, ":")))

	err = m.db.Update(func(txn *badger.Txn) error {
		return txn.SetEntry(&badger.Entry{
			Key:       []byte(tokenString),
			Value:     []byte(value),
			ExpiresAt: uint64(time.Now().Add(dur).Unix()),
		})
	})

	return tokenString, err
}

// Validate implements service.TokenStore.
func (m *Manager) Validate(ctx context.Context, tokenString string) (userID string, filename string, ok bool) {

	err := m.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(tokenString))
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			value, err := base64.URLEncoding.DecodeString(string(val))
			if err != nil {
				return err
			}
			res := strings.Split(string(value), ":")
			userID = res[0]
			filename = res[1]
			ok = true

			return nil
		})
	})

	if err != nil {
		if err != badger.ErrKeyNotFound {
			panic(err)
		}
	}
	return
}
