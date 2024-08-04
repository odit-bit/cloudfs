package repo

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/glebarez/go-sqlite"
	"github.com/odit-bit/cloudfs/internal/user"
)

type DB struct {
	sqlite *sql.DB
}

func NewSQLiteDB(path string) (DB, error) {
	db, err := DefaultDB("rwc", path)
	if err != nil {
		return DB{}, err
	}

	us, err := newAccountDB(db)
	if err != nil {
		return DB{}, err
	}
	return us, nil
}

func newAccountDB(db *sql.DB) (DB, error) {
	query := `	
	CREATE TABLE IF NOT EXISTS Account (
			ID TEXT PRIMARY KEY NOT NULL UNIQUE,
			Name TEXT UNIQUE,
			HashPassword BLOB
		);
	`

	if _, err := db.ExecContext(context.Background(), query); err != nil {
		return DB{}, err
	}

	adb := DB{
		sqlite: db,
	}

	return adb, nil
}

func DefaultDB(mode, path string) (*sql.DB, error) {
	dsn := ":memory:"
	if path != "" {
		dsn = fmt.Sprintf("file:%s?cache=shared&mode=%v", path, mode)
	}

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (db *DB) Find(ctx context.Context, name string) (*user.Account, error) {
	var account user.Account

	row := db.sqlite.QueryRow("SELECT * FROM Account WHERE Name = ? LIMIT 1", name)
	err := row.Scan(&account.ID, &account.Name, &account.HashPassword)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (db *DB) Insert(ctx context.Context, account *user.Account) error {
	query := "INSERT INTO Account (ID, Name, HashPassword) VALUES (?, ?, ?)"
	_, err := db.sqlite.ExecContext(ctx, query, account.ID, account.Name, account.HashPassword)
	if err != nil {
		return err
	}

	return nil
}
