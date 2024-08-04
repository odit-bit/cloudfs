package repo

import (
	"context"
	"testing"
	"time"

	"github.com/odit-bit/cloudfs/internal/token"
	"github.com/stretchr/testify/assert"
)

func Test_token(t *testing.T) {
	m, err := NewInMemToken("")
	if err != nil {
		t.Fatal(err)
	}
	defer m.Close()

	tkn := token.NewShareToken("odit", "filename", time.Duration(10*time.Second))
	if err := m.Put(context.Background(), tkn); err != nil {
		t.Fatal(err)
	}

	actual, _, err := m.Get(context.Background(), tkn.Key())
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, tkn.ValidUntil().Unix(), actual.ValidUntil().Unix())
}
