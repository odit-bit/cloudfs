package repo

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/odit-bit/cloudfs/internal/user"
)

const default_table_name = "accounts"

type userPG struct {
	pg *sql.DB
}

func NewUserPG(ctx context.Context, uri string) (*userPG, error) {
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

	adb := userPG{
		pg: db,
	}

	return &adb, nil
}

func migratePG(db *sql.DB, tableName string) error {
	extentionSTMT := `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`
	if _, err := db.ExecContext(context.Background(), extentionSTMT); err != nil {
		return err
	}

	// query := `
	// CREATE TABLE IF NOT EXISTS account (
	// 		ID UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
	// 		Name VARCHAR (50) UNIQUE NOT NULL,
	// 		HashPassword VARCHAR NOT NULL
	// 	);
	// `
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %v (
		ID uuid PRIMARY KEY DEFAULT uuid_generate_v4 (),
		Name VARCHAR (50) UNIQUE NOT NULL,
		HashPassword VARCHAR NOT NULL
	);`, tableName)

	if _, err := db.ExecContext(context.Background(), query); err != nil {
		return err
	}

	return nil
}

func (db *userPG) Close() error {
	return db.pg.Close()
}

func (db *userPG) Find(ctx context.Context, name string) (*user.Account, error) {
	var account user.Account

	row := db.pg.QueryRow("SELECT * FROM accounts WHERE Name = $1 LIMIT 1", name)
	err := row.Scan(&account.ID, &account.Name, &account.HashPassword)
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (db *userPG) Insert(ctx context.Context, account *user.Account) error {
	query := "INSERT INTO accounts (ID, Name, HashPassword) VALUES ($1, $2, $3)"
	_, err := db.pg.ExecContext(ctx, query, account.ID, account.Name, account.HashPassword)
	if err != nil {
		return err
	}

	return nil
}
