package repo

import (
	"context"
	"database/sql"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/odit-bit/cloudfs/internal/user"
	"github.com/stretchr/testify/assert"
)

func Test_token(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbUrl := "postgres://admin:admin@localhost:5432/postgres?sslmode=disable"
	db, err := sql.Open("pgx", dbUrl)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	defer dropTable(db, "user_tokens")
	pg, err := NewUserTokenPG(ctx, db)
	if err != nil {
		t.Fatal(err)
	}

	expect := user.NewToken("id-124", 1*time.Second)
	if err := pg.PutToken(ctx, expect); err != nil {
		t.Fatal(err)
	}

	actual, err := pg.GetToken(ctx, expect.Key)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, expect, actual)

	actual2, ok := pg.GetTokenUserID(ctx, expect.UserID)
	if !ok {
		t.Fatal("should ok")
	}
	assert.Equal(t, expect, actual2)

	if err := pg.Delete(ctx, expect.Key); err != nil {
		t.Fatal(err)
	}

}
