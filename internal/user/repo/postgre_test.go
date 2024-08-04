package repo

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/odit-bit/cloudfs/internal/user"
)

const default_host = "localhost"
const default_port = 5432

func dropTable(db *sql.DB, tableName string) error {
	query := fmt.Sprintf("DROP TABLE IF EXISTS %v;", tableName)
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func Test_pguser(t *testing.T) {
	dbUrl := "postgres://admin:admin@localhost:5432/postgres"
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	userDB, _ := NewUserPG(ctx, dbUrl)
	acc1 := user.CreateAccount("user1", "12345")
	if err := userDB.Insert(ctx, acc1); err != nil {
		t.Fatal(err)
	}

	acc2, err := userDB.Find(ctx, "user1")
	if err != nil {
		t.Fatal(err)
	}

	if ok := acc2.CheckPassword("12345"); !ok {
		t.Fatal("wrong password")
	}

}
