package user

import (
	"context"
	"database/sql"
	"fmt"
)

type DB struct {
	sqlite *sql.DB
}

func NewDB() (*DB, error) {
	db, err := DefaultDB("rwc")
	if err != nil {
		return nil, err
	}

	us, err := newAccountDB(db)
	if err != nil {
		return nil, err
	}
	return us, nil
}

func newAccountDB(db *sql.DB) (*DB, error) {
	query := `	
	CREATE TABLE IF NOT EXISTS Account (
			ID TEXT PRIMARY KEY NOT NULL UNIQUE,
			Name TEXT UNIQUE,
			HashPassword BLOB
		);
	`

	if _, err := db.ExecContext(context.Background(), query); err != nil {
		return nil, err
	}

	adb := DB{
		sqlite: db,
	}

	return &adb, nil
}

func DefaultDB(mode string) (*sql.DB, error) {
	defaultPath := fmt.Sprintf("file:account.db?cache=shared&mode=%v", mode)

	// Open a test database (this creates a new database in memory)
	db, err := sql.Open("sqlite3", defaultPath)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (db *DB) Find(ctx context.Context, name string) (*Account, error) {
	var account Account

	row := db.sqlite.QueryRow("SELECT * FROM Account WHERE Name = ? LIMIT 1", name)
	err := row.Scan(&account.ID, &account.Name, &account.HashPassword)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (db *DB) Insert(ctx context.Context, account *Account) error {
	query := "INSERT INTO Account (ID, Name, HashPassword) VALUES (?, ?, ?)"
	_, err := db.sqlite.ExecContext(ctx, query, account.ID, account.Name, account.HashPassword)
	if err != nil {
		return err
	}

	return nil
}
