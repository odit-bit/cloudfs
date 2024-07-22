package tokenpg

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/odit-bit/cloudfs/internal/blob"
	"github.com/odit-bit/cloudfs/service"
)

const default_table_name = "share_tokens"

var _ service.TokenStore = (*DB)(nil)

type DB struct {
	*sql.DB
}

func NewDB(ctx context.Context, uri string) (*DB, error) {
	db, err := sql.Open("pgx", uri)
	if err != nil {
		return nil, err
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	if err := migratePG(db, default_table_name); err != nil {
		return nil, err
	}

	adb := DB{
		DB: db,
	}

	return &adb, nil
}

func migratePG(db *sql.DB, tableName string) error {
	// extentionSTMT := `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`
	// if _, err := db.ExecContext(context.Background(), extentionSTMT); err != nil {
	// 	return err
	// }

	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %v (
		key VARCHAR (50) PRIMARY KEY UNIQUE NOT NULL,
		user_id VARCHAR(50) NOT NULL,
		filename VARCHAR(128) UNIQUE NOT NULL,
		expire_at INT NOT NULL

	);`, tableName)

	if _, err := db.ExecContext(context.Background(), query); err != nil {
		return err
	}

	return nil
}

type txn struct {
	tx *sql.Tx
}

func (txn *txn) Commit() error {
	return txn.tx.Commit()
}
func (txn *txn) Cancel() error {
	return txn.tx.Rollback()
}

func (txn *txn) Put(ctx context.Context, token *blob.ShareToken) error {
	query := `
		INSERT INTO share_tokens (key, user_id, filename, expire_at) 
		VALUES ($1, $2, $3, $4)
		ON CONFLICT(filename)
		DO UPDATE SET 
			key = EXCLUDED.key,
			expire_at = EXCLUDED.expire_at;
	`
	_, err := txn.tx.ExecContext(ctx, query, token.Key, token.UserID, token.Filename, token.Expire.Unix())
	if err != nil {
		return err
	}
	return nil
}

func (txn *txn) Get(ctx context.Context, tokenString string) (*blob.ShareToken, error) {

	row := txn.tx.QueryRowContext(ctx, "SELECT * FROM share_tokens WHERE key = $1 LIMIT 1", tokenString)
	var st blob.ShareToken
	var unix int64
	err := row.Scan(&st.Key, &st.UserID, &st.Filename, &unix)
	if err != nil {
		return nil, err
	}
	st.Expire = time.Unix(unix, 0)

	return &st, nil
}

func (txn *txn) Delete(ctx context.Context, key string) error {
	_, err := txn.tx.ExecContext(ctx, "DELETE FROM share_tokens WHERE key = $1;")
	return err
}

func (t *DB) Query(ctx context.Context, fn func(txn service.TokenTxn) error) error {

	tx, err := t.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	txn := txn{
		tx: tx,
	}

	return fn(&txn)

}

// Generate implements service.TokenStore.
func (t *DB) Generate(ctx context.Context, bucket string, filename string, dur time.Duration) (string, error) {
	panic("unimplemented")
}

// Validate implements service.TokenStore.
func (t *DB) Validate(ctx context.Context, tokenString string) (userID string, filename string, ok bool) {
	panic("unimplemented")
}
