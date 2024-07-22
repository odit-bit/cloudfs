package tokenbunt

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_token(t *testing.T) {
	m, err := New("")
	if err != nil {
		t.Fatal(err)
	}
	defer m.db.Close()

	bucket := "my-bucket"
	filename := "my-file-1"
	tkn, err := m.Generate(context.Background(), bucket, filename, 1*time.Minute)
	if err != nil {
		t.Fatal(err)
	}

	actID, actFilename, ok := m.Validate(context.Background(), tkn)
	assert.Equal(t, true, ok)
	assert.Equal(t, filename, actFilename)
	assert.Equal(t, bucket, actID)

	tkn, err = m.Generate(context.Background(), bucket, filename, 1*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(2 * time.Millisecond)
	_, _, ok = m.Validate(context.Background(), tkn)
	assert.Equal(t, false, ok)
}
