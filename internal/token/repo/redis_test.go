package repo

import (
	"context"
	"testing"
	"time"

	"github.com/odit-bit/cloudfs/internal/token"
	"github.com/stretchr/testify/assert"
)

func Test_tokenRedis(t *testing.T) {
	uri := "redis://:@localhost:6379/0" // password set
	tr := NewRedisToken(uri)

	tkn := token.NewShareToken("123-d", "file", 1*time.Minute)
	if err := tr.Put(context.Background(), tkn); err != nil {
		t.Fatal(err)
	}
	key := tkn.Key()
	actual, _, err := tr.Get(context.TODO(), key)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, tkn.Filename(), actual.Filename())
}
