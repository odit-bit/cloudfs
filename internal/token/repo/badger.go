package repo

import (
	"bytes"
	"context"
	"encoding/gob"
	"log"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/odit-bit/cloudfs/internal/token"
)

// var _ service.TokenStore = (*Manager)(nil)
// var _ service.TokenTxn = (*Manager)(nil)

type badgerToken struct {
	Key      string
	UserID   string
	Filename string
	Expire   time.Time
	value    []byte
}

func toBadger(tkn *token.ShareToken) (*badgerToken, error) {
	bt := badgerToken{
		Key:      tkn.Key(),
		UserID:   tkn.UserID(),
		Filename: tkn.Filename(),
		Expire:   tkn.ValidUntil(),
	}

	buf := bytes.Buffer{}
	if err := gob.NewEncoder(&buf).Encode(bt); err != nil {
		return nil, err
	}
	bt.value = buf.Bytes()
	return &bt, nil
}

func (bt *badgerToken) keyBytes() []byte {
	return []byte(bt.Key)
}

func (bt *badgerToken) valueBytes() []byte {
	return bt.value
}

type Manager struct {
	db *badger.DB
}

// // Cancel implements service.TokenTxn.
// func (m *Manager) Cancel() error {
// 	panic("unimplemented")
// }

// // Commit implements service.TokenTxn.
// func (m *Manager) Commit() error {
// 	panic("unimplemented")
// }

// // Delete implements service.TokenTxn.
// func (m *Manager) Delete(ctx context.Context, key string) error {
// 	panic("unimplemented")
// }

// Get implements service.TokenTxn.
func (m *Manager) Get(ctx context.Context, tokenString string) (*token.ShareToken, bool, error) {
	var ok bool
	tkn := new(token.ShareToken)
	err := m.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(tokenString))
		if err != nil {
			if err == badger.ErrKeyNotFound {
				return nil
			}
			return err
		}

		// value
		if err := item.Value(func(val []byte) error {
			bt := badgerToken{}
			if err := gob.NewDecoder(bytes.NewReader(val)).Decode(&bt); err != nil {
				return err
			}
			tkn = token.FromStore(bt.Key, bt.UserID, bt.Filename, bt.Expire)
			return nil

		}); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, ok, err
	}
	return tkn, ok, nil
}

// Put implements service.TokenTxn.
func (m *Manager) Put(ctx context.Context, token *token.ShareToken) error {

	return m.db.Update(func(txn *badger.Txn) error {
		bt, err := toBadger(token)
		if err != nil {
			return err
		}
		return txn.SetEntry(&badger.Entry{
			Key:       bt.keyBytes(),
			Value:     bt.valueBytes(),
			ExpiresAt: uint64(bt.Expire.Unix()),
			UserMeta:  0,
		})

	})
}

// // Query implements service.TokenStore.
// func (m *Manager) Query(ctx context.Context, txn func(txn service.TokenTxn) error) error {
// 	panic("unimplemented")
// }

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

func NewInMemToken(path string) (Manager, error) {
	return Manager{
		db: connectBadger(path),
	}, nil
}

// // Generate implements service.TokenStore.
// func (m *Manager) Generate(ctx context.Context, bucket string, filename string, dur time.Duration) (string, error) {

// 	//make token from this enc
// 	b := make([]byte, 16)
// 	_, err := rand.Read(b)
// 	if err != nil {
// 		return "", err
// 	}
// 	tokenString := hex.EncodeToString(b)
// 	value := base64.URLEncoding.EncodeToString([]byte(strings.Join([]string{bucket, filename}, ":")))

// 	err = m.db.Update(func(txn *badger.Txn) error {
// 		return txn.SetEntry(&badger.Entry{
// 			Key:       []byte(tokenString),
// 			Value:     []byte(value),
// 			ExpiresAt: uint64(time.Now().Add(dur).Unix()),
// 		})
// 	})

// 	return tokenString, err
// }

// // Validate implements service.TokenStore.
// func (m *Manager) Validate(ctx context.Context, tokenString string) (userID string, filename string, ok bool) {

// 	err := m.db.View(func(txn *badger.Txn) error {
// 		item, err := txn.Get([]byte(tokenString))
// 		if err != nil {
// 			return err
// 		}

// 		return item.Value(func(val []byte) error {
// 			value, err := base64.URLEncoding.DecodeString(string(val))
// 			if err != nil {
// 				return err
// 			}
// 			res := strings.Split(string(value), ":")
// 			userID = res[0]
// 			filename = res[1]
// 			ok = true

// 			return nil
// 		})
// 	})

// 	if err != nil {
// 		if err != badger.ErrKeyNotFound {
// 			panic(err)
// 		}
// 	}
// 	return
// }
