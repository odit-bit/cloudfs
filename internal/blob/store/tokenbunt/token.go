package tokenbunt

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"strings"
	"time"

	"github.com/odit-bit/cloudfs/service"
	"github.com/tidwall/buntdb"
)

var _ service.TokenStore = (*Manager)(nil)

type Manager struct {
	db *buntdb.DB
}

func (m *Manager) Close() error {
	return m.db.Close()
}

// Generate implements service.TokenStore.
func (m *Manager) Generate(ctx context.Context, bucket string, filename string, dur time.Duration) (string, error) {
	value := base64.URLEncoding.EncodeToString([]byte(strings.Join([]string{bucket, filename}, ":")))
	//make token from this enc
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	key := hex.EncodeToString(b)

	err = m.db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(key, value, &buntdb.SetOptions{
			Expires: true,
			TTL:     dur,
		})
		return err
	})

	return key, err
}

func New(path string) (*Manager, error) {
	if path == "" {
		path = ":memory:"
	}
	bdb, err := buntdb.Open(path)
	if err != nil {
		return nil, err
	}
	return &Manager{
		db: bdb,
	}, nil
}

// Validate implements service.TokenStore.
func (m *Manager) Validate(ctx context.Context, tokenString string) (userID string, filename string, ok bool) {
	m.db.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get(tokenString, false)
		if err != nil {
			return err
		}

		value, err := base64.URLEncoding.DecodeString(val)
		if err != nil {
			return err
		}
		res := strings.Split(string(value), ":")
		userID = res[0]
		filename = res[1]
		ok = true

		return nil
	})
	return
}
